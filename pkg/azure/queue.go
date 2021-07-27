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

	response, err := queueURL.Create(getContext(), getQueueMetadata())
	if err != nil {
		if azErr, ok := err.(azqueue.StorageError); ok {
			if azErr.ServiceCode() == azqueue.ServiceCodeQueueAlreadyExists  {
				fmt.Println("Queue already exists")
				return nil
			}
			return azErr
		}
	}
	fmt.Println("Queue created: ", response.RequestID())
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

	response, err := messagesURL.Enqueue(getContext(), message, 0, 60 * time.Second)

	if err != nil {
		return err
	}

	fmt.Println("Message posted: ", response.RequestID())
	return nil
}