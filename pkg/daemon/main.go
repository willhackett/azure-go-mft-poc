package daemon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/sirupsen/logrus"
	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
	"github.com/willhackett/azure-mft/pkg/logger"
)

func canAgentSendFile(agentName string) bool {
	cfg := config.GetConfig()

	return constant.StringInList(agentName, cfg.AllowRequestsFrom)
}

func canAgentRequestFile(agentName string) bool {
	cfg := config.GetConfig()

	return constant.StringInList(agentName, cfg.AllowRequestsFrom) || agentName == cfg.Agent.Name
}

func handleMessage(qm *QueueMessage) {
	log := logger.Get().WithFields(logrus.Fields{
		"event": "HandleMessage",
	})

	messageBody := constant.Message{}

	// Unmarshal message from JSON
	err := json.Unmarshal([]byte(qm.text), &messageBody)
	if err != nil {
		log.Trace(err)
		log.WithField("message_body", qm.text).Warn("Discarded invalid message payload")
		return
	}

	// Verify the contents of the message signature
	err = keys.VerifyMessage(messageBody)
	if err != nil {
		log.Trace(err)
		log.WithField("id", messageBody.ID).WithField("body", qm.text).Warn("Message signature cannot be verified")
		return
	}

	switch messageBody.Type {
	case constant.FileRequestMessageType:
		// Check if requesting agent is allowed to request files
		if !canAgentRequestFile(messageBody.Agent) {
			log.WithField("id", messageBody.ID).WithField("destination_agent", messageBody.Agent).Warn("Requesting agent is not allowed to request files")
			return
		}

		err = handleFileRequest(messageBody)
	case constant.FileHandshakeMessageType:
		// Check if requesting agent is allowed to send files
		if !canAgentSendFile(messageBody.Agent) {
			log.WithField("id", messageBody.ID).WithField("destination_agent", messageBody.Agent).Warn("Requesting agent is not allowed to send files")
			return
		}

		err = handleFileHandshake(messageBody)
	case constant.FileHandshakeResponseMessageType:
		handleFileHandshakeResponse(qm, messageBody)

	case constant.FileAvailableMessageType:
		if !canAgentSendFile(messageBody.Agent) {
			log.WithField("id", messageBody.ID).WithField("destination_agent", messageBody.Agent).Warn("Requesting agent is not allowed to request files")
			return
		}

		err = handleFileAvailable(qm, messageBody)
	default:
		log.WithField("id", messageBody.ID).WithField("body", qm.text).Warn("Invalid Type on Message")
		return
	}

	if err != nil {
		log.WithField("id", messageBody.ID).Warn("Failed to process message, releasing to queue")
		return
	}

	log.WithField("id", messageBody.ID).Info("Successful " + messageBody.Type + " operation")
	log.WithField("id", messageBody.ID).Debug("Discarding dequeued message")
	qm.URL.Delete(qm.context, qm.popReceipt)
}

func Init() {
	messagesURL, azureContext := azure.GetMessagesURLAndContext()

	messageChannel := make(chan *azqueue.DequeuedMessage, constant.MaxConcurrentTransfers)

	log := logger.Get().WithFields(logrus.Fields{
		"event": "QueueOperation",
	})

	for i := 0; i < constant.MaxConcurrentTransfers; i++ {
		// Go routine for handling messages
		go func(messageChannel <-chan *azqueue.DequeuedMessage) {
			for {
				inboundMessage := <-messageChannel
				popReceipt := inboundMessage.PopReceipt
				URL := messagesURL.NewMessageIDURL(inboundMessage.ID)

				queueMessage := &QueueMessage{
					context:    azureContext,
					text:       inboundMessage.Text,
					popReceipt: popReceipt,
					URL:        URL,
				}

				if inboundMessage.DequeueCount > constant.MaxRetriesThreshold {
					URL.Delete(azureContext, popReceipt)
					log.Warn(fmt.Sprintf("Deleted message with ID: %s as it reached the failure threshold %d", inboundMessage.ID, constant.MaxRetriesThreshold))
					continue
				}

				handleMessage(queueMessage)
			}
		}(messageChannel)
	}

	for {
		// Try to dequeue a batch of messages from the queue
		dequeue, err := messagesURL.Dequeue(azureContext, azqueue.QueueMaxMessagesDequeue, 60*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		if dequeue.NumMessages() == 0 {
			time.Sleep(time.Second * 1)
		} else {
			log.Debug("Processing new messages")

			for m := int32(0); m < dequeue.NumMessages(); m++ {
				messageChannel <- dequeue.Message(m)
			}
		}
	}
}
