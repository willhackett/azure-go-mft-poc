package azure

import (
	"fmt"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
)

func getBlobMetadata() azblob.Metadata {
	return azblob.Metadata{
		"agent": config.GetConfig().Agent.Name,
	}
}

func getContainer(containerName string) *azblob.ContainerURL {
	URL := getContainerURL(containerName)

	container := azblob.NewContainerURL(URL, azurePipeline)

	return &container
}

func getBlobURL(containerName string, blobName string) azblob.BlobURL {
	container := getContainer(containerName)

	return container.NewBlobURL(blobName)
}

// UpsertContainer creates a blob storage container if it does not already exist
func UpsertContainer(containerName string) error {
	container := getContainer(containerName)

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
	err := getCredential()
	if err != nil {
		cobra.CheckErr(err)
	}

	azurePipeline = azblob.NewPipeline(azureCredential, azblob.PipelineOptions{})

	// Create agent container
	if err := UpsertContainer(config.GetConfig().Agent.Name); err != nil {
		cobra.CheckErr(err)
	}
	// Create public keys container
	if err := UpsertContainer("publickeys"); err != nil {
		cobra.CheckErr(err)
	}
}

// func UploadFile(containerName string, blobName string) error {
// 	blobURL := getBlobURL(containerName, blobName)

// }

func UploadBuffer(containerName string, blobName string, buffer []byte) error {
	blockBlobURL := getBlobURL(containerName, blobName).ToBlockBlobURL()

	response, err := azblob.UploadBufferToBlockBlob(getContext(), buffer, blockBlobURL, azblob.UploadToBlockBlobOptions{
		BlockSize: 2 * 1024,
		Metadata:  getBlobMetadata(),
	})

	if err != nil && response.RequestID() != "" {
		return err
	}

	return nil
}

func DownloadBuffer(containerName string, blobName string) ([]byte, error) {
	blobURL := getBlobURL(containerName, blobName)

	properties, err := blobURL.GetProperties(azureContext, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, properties.ContentLength())
	err = azblob.DownloadBlobToBuffer(getContext(), blobURL, 0, 0, bytes, azblob.DownloadFromBlobOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bytes, nil
}

func UploadFromFile(containerName string, blobName string, fileName string) error {
	blockBlobURL := getBlobURL(containerName, blobName).ToBlockBlobURL()

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	response, err := azblob.UploadFileToBlockBlob(getContext(), file, blockBlobURL, azblob.UploadToBlockBlobOptions{
		BlockSize: 32 * 1024,
		Metadata:  getBlobMetadata(),
		Progress: func(bytes int64) {
			fmt.Printf("\rUploading: %d", bytes)
		},
	})

	sasToken, err := getBlobURL(containerName, blobName)

	if err != nil {
		return err
	}

	fmt.Println("\nUpload complete: ", response.RequestID())

	return nil
}
