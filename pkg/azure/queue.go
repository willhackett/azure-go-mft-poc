package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
)

func getQueueMetadata() azqueue.Metadata {
	return azqueue.Metadata{
		"agent": config.GetConfig().Agent.Name,
	}
}

func getQueue(queueName string) azqueue.QueueURL {
	return azqueue.NewQueueURL(getQueueURL(queueName), azurePipeline)
}

func getMessagesURL(queueName string) azqueue.MessagesURL {
	return getQueue(queueName).NewMessagesURL()
}

// Upsert creates a queue if it does not exist
func UpsertQueue(queueName string) error {
	queueURL := getQueue(queueName)

	_, err := queueURL.Create(getContext(), getQueueMetadata())
	if err != nil {
		if azErr, ok := err.(azqueue.StorageError); ok {
			if azErr.ServiceCode() == azqueue.ServiceCodeQueueAlreadyExists {
				log.Debug(fmt.Sprintf("Queue %s already exists", queueName))
				return nil
			}
			return azErr
		}
	}
	log.Debug(fmt.Sprintf("Queue %s created", queueName))
	return nil
}

// InitQueue is called by Cobra to setup the queue as needed
func InitQueue() {
	if err := UpsertQueue(config.GetConfig().Agent.Name); err != nil {
		cobra.CheckErr(err)
	}
}

func GetMessagesURLAndContext() (azqueue.MessagesURL, context.Context) {
	messagesURL := getMessagesURL(config.GetConfig().Agent.Name)

	return messagesURL, azureContext
}

func PostMessage(queueName string, message string) error {
	messagesURL := getMessagesURL(queueName)

	_, err := messagesURL.Enqueue(getContext(), message, 0, 60*time.Minute)

	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Message %s posted to queue %s", message, queueName))
	return nil
}
