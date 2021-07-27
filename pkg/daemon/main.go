package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
)

func handleFileRequest() {

}

func handleFileHandshake() {

}

func handleFileHandshakeResponse() {

}

func handleFileAvailable() {

}

func canAgentSendFile(agentName string) bool {
	cfg := config.GetConfig()

	return constant.StringInList(agentName, cfg.AllowRequestsFrom)
}

func canAgentRequestFile(agentName string) bool {
	cfg := config.GetConfig()

	return constant.StringInList(agentName, cfg.AllowRequestsFrom) || agentName == cfg.Agent.Name
}

func handleMessage(message *Message) {
	messageBody := constant.Message{}
	err := json.Unmarshal([]byte(message.text), &messageBody)
	if err != nil {
		fmt.Println("Invalid message body, discarding")
	}

	err = keys.VerifyMessage(messageBody)

	if err != nil {
		fmt.Println("Message signature not accepted, rejecting")
		return
	}

	switch messageBody.Type {
		case constant.FileRequestMessageType:
			// Check if requesting agent is allowed to request files
			if !canAgentRequestFile(messageBody.Agent) {
				fmt.Println("Agent not allowed to request files, rejecting")
				return
			}

			handleFileRequest()
		case constant.FileHandshakeMessageType:
			// Check if requesting agent is allowed to send files
			if !canAgentSendFile(messageBody.Agent) {
				fmt.Println("Agent not allowed to request files, rejecting")
				return
			}
			
			handleFileHandshake()
		case constant.FileHandshakeResponseMessageType:
			handleFileHandshakeResponse()

		case constant.FileAvailableMessageType:
			if !canAgentSendFile(messageBody.Agent) {
				fmt.Println("Agent not allowed to request files, rejecting")
				return
			}

			handleFileAvailable()
		default:
			fmt.Println("Invalid message type, discarding")
	}
}

func Init() {
	messagesURL, azureContext := azure.GetMessagesURLAndContext()
	
	messageChannel := make(chan *azqueue.DequeuedMessage, constant.MaxConcurrentTransfers)

	for i := 0; i < constant.MaxConcurrentTransfers; i++ {
		// Go routine for handling messages
		go func(messageChannel <-chan *azqueue.DequeuedMessage) {
			for {
				inboundMessage := <-messageChannel
				popReceipt := inboundMessage.PopReceipt
				URL := messagesURL.NewMessageIDURL(inboundMessage.ID)

				message := &Message{
					context: azureContext,
					text: inboundMessage.Text,
					popReceipt: &popReceipt,
					URL: URL,
				}

				if inboundMessage.DequeueCount > constant.MaxRetriesThreshold {
					URL.Delete(azureContext, popReceipt)
					fmt.Println("Deleted " + inboundMessage.ID + " because it reached the failure theashold")
					continue
				}

				handleMessage(message)
			}
		}(messageChannel)
	}

	for {
		// Try to dequeue a batch of messages from the queue
		dequeue, err := messagesURL.Dequeue(azureContext, azqueue.QueueMaxMessagesDequeue, 10*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		if dequeue.NumMessages() == 0 {
			fmt.Println("Queue empty, waiting 10 seconds")
			time.Sleep(time.Second * 10)
		} else {
			fmt.Println("messages on queue, pulling")

			for m := int32(0); m < dequeue.NumMessages(); m++ {
				messageChannel <- dequeue.Message(m)
			}
		}
	}
}