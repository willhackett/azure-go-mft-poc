package azure

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/willhackett/azure-mft/pkg/config"
)

var (
	azureCredential *azblob.SharedKeyCredential
	azurePipeline pipeline.Pipeline
	azureContext context.Context
)

func getCredential() error {
	accountName, accountKey := config.GetConfig().Azure.AccountName, config.GetConfig().Azure.AccountKey

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}

	azureCredential = credential

	return nil
}

func getResourceURL(name string, resource string) url.URL {
	accountName := config.GetConfig().Azure.AccountName

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.%s.core.windows.net/%s", accountName, resource, name),
	)

	return *URL
}

func getContainerURL(containerName string) url.URL {
	return getResourceURL(containerName, "blob")
}

func getQueueURL(queueName string) url.URL {
	return getResourceURL(queueName, "queue")
}

func getContext() context.Context {
	if azureContext != nil {
		return azureContext
	}

	azureContext = context.Background()

	return azureContext
}
