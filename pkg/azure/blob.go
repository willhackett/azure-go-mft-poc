package azure

import (
	"fmt"
	"net/url"
	"os"
	"time"

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

func UploadFromFile(containerName string, blobName string, fileName string, progress func(bytes int64)) (string, error) {
	blobURL := getBlobURL(containerName, blobName)
	blockBlobURL := blobURL.ToBlockBlobURL()

	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	response, err := azblob.UploadFileToBlockBlob(getContext(), file, blockBlobURL, azblob.UploadToBlockBlobOptions{
		BlockSize: 32 * 1024,
		Metadata:  getBlobMetadata(),
		Progress:  progress,
	})
	if err != nil {
		fmt.Println("Failed uploading file", err)
		return "", err
	}

	file.Close()

	fmt.Println("Successfully uploaded blob", response.RequestID())

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		Permissions:   azblob.BlobSASPermissions{Read: true}.String(),
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(azureCredential)
	if err != nil {
		fmt.Println("Failed generating SAS QPs", err)
		return "", err
	}

	signedURL := sasQueryParams.Encode()

	signedURL = fmt.Sprintf(
		"https://%s.blob.core.windows.net/%s/%s?%s",
		config.GetConfig().Azure.AccountName,
		containerName,
		blobName,
		signedURL,
	)

	if err != nil {
		return "", err
	}

	fmt.Println("\nUpload complete: ", response.RequestID())

	return signedURL, nil
}

func DownloadSignedURLToFile(signedURL string, fileName string, progress func(bytes int64)) error {
	blobURLBase, _ := url.Parse(signedURL)
	anonymousCredential := azblob.NewAnonymousCredential()
	pipeline := azblob.NewPipeline(anonymousCredential, azblob.PipelineOptions{})
	blobURL := azblob.NewBlobURL(*blobURLBase, pipeline)

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed creating file", err)
		return nil
	}

	defer file.Close()

	err = azblob.DownloadBlobToFile(azureContext, blobURL, 0, 0, file, azblob.DownloadFromBlobOptions{
		BlockSize: 32 * 1024,
		Progress:  progress,
	})

	if err != nil {
		fmt.Println("Failed downloading file", err)
		return nil
	}

	fmt.Println("Successfully downloaded blob")
	return nil
}
