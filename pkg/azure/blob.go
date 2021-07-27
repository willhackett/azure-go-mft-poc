package azure

import (
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/willhackett/azure-mft/pkg/config"
)

func getBlobMetadata() azblob.Metadata {
	return azblob.Metadata{
		"agent": config.GetConfig().Agent.Name,
	}
}

// UpsertContainer creates a blob storage container if it does not already exist
func UpsertContainer(containerName string) error {
	pipeline, err := getPipeline()
	if err != nil {
		return err
	}

	URL := getContainerURL(containerName)
	
	container := azblob.NewContainerURL(URL, pipeline)

	response, err := container.Create(getContext(), getBlobMetadata(), azblob.PublicAccessNone)

	if err != nil {
		if azErr, ok := err.(azblob.StorageError); ok {
			if azErr.ServiceCode() == azblob.ServiceCodeContainerAlreadyExists {
				fmt.Println("Storage container already exists")
				return nil
			}
			return azErr
		}
	}

	fmt.Println("Storage container created: ", response.RequestID())
	return nil
}

// InitBlob creates the blob containers if they do not exist
func InitBlob() {
	// Create agent container
	if err := UpsertContainer(config.GetConfig().Agent.Name); err != nil {
		fmt.Println("Error creating Azure container: ", err)
	}
	// Create public keys container
	if err := UpsertContainer("publickeys"); err != nil {
		fmt.Println("Error creating Azure container: ", err)
	}
}