package azure

import (
	"fmt"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
)

func getQueueMetadata() azqueue.Metadata {
	return azqueue.Metadata{
		"agent": config.GetConfig().Agent.Name,
	}
}

// Upsert creates a queue if it does not exist
func UpsertQueue(queueName string) error {
	queueURL := azqueue.NewQueueURL(getQueueURL(queueName), azurePipeline)

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