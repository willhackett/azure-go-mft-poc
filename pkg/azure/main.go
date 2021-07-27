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

func getCredential() (*azblob.SharedKeyCredential, error) {
	if azureCredential != nil {
		return azureCredential, nil
	}

	accountName, accountKey := config.GetConfig().Azure.AccountName, config.GetConfig().Azure.AccountKey

	azureCredential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	return azureCredential, nil
}

func getPipeline() (pipeline.Pipeline, error) {
	if azurePipeline != nil {
		return azurePipeline, nil
	}

	credential, err := getCredential()
	if err != nil {
		return nil, err
	}

	azurePipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{})
	return azurePipeline, nil
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
