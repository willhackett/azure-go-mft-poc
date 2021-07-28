package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
	"github.com/willhackett/azure-mft/pkg/registry"
	"github.com/willhackett/azure-mft/pkg/tasks"
)

func handleFileRequest(m constant.Message) error {
	body := constant.FileRequestMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

	registry.AddTransfer(m.ID, body, 5*60*60*1000)

	fmt.Println("Received file send request", body)

	file, err := os.Open(body.FileName)
	if err != nil {
		fmt.Println("Cannot open file", err)
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Cannot get file info", err)
		return err
	}

	fileSize := fileInfo.Size()

	err = tasks.SendFileHandshake(m.ID, body.DestinationFileName, fileSize, body.DestinationAgent)
	if err != nil {
		fmt.Println("Cannot send file handshake", err)
		return err
	}

	return nil
}

func handleFileHandshake(m constant.Message) error {
	body := constant.FileHandshakeMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}
	fmt.Println("File Handshake", body)

	_, err := os.OpenFile(body.FileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("ERROR", err)
		tasks.SendFileHandshakeResponse(m.ID, false, m.Agent, fmt.Sprintf("Cannot open destination path: %s", err))
		return nil
	}

	if err = tasks.SendFileHandshakeResponse(m.ID, true, m.Agent, ""); err != nil {
		return err
	}

	return nil
}

func handleFileHandshakeResponse(qm *QueueMessage, m constant.Message) error {
	body := constant.FileHandshakeResponseMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}
	fmt.Println("File Handshake Response", body)

	if !body.Accepted {
		fmt.Println("Rejected, transfer to be considered failed.")
		return nil
	}

	transfer, ok := registry.GetTransfer(m.ID)
	if !ok {
		fmt.Println("Cannot get details of transfer, perhaps it expired")
		return errors.New("transfer expired or did not originate from this node")
	}

	debounce := time.Now().Add(time.Second * 30).Unix()

	reportProgress := func(bytes int64) {
		if time.Now().Unix() > debounce {
			fmt.Println("Upload bytes", bytes)
			qm.IncreaseLease()
			debounce = time.Now().Add(time.Second * 30).Unix()
		}
	}

	signedURL, err := azure.UploadFromFile(transfer.Details.DestinationAgent, m.ID, transfer.Details.FileName, reportProgress)
	if err != nil {
		fmt.Println("Failed to upload", err)
		return err
	}

	encryptedSignedURL, err := keys.EncryptString(transfer.Details.DestinationAgent, m.KeyID, signedURL)
	if err != nil {
		fmt.Println("FAILED TO ENCRYPT", err)
		return err
	}

	err = tasks.SendFileAvailable(m.ID, encryptedSignedURL, transfer.Details.DestinationFileName, transfer.Details.DestinationAgent)

	return err
}

func handleFileAvailable(qm *QueueMessage, m constant.Message) error {
	body := constant.FileAvailableMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

	signedURL, err := keys.DecryptString(body.SignedURL)
	if err != nil {
		fmt.Println("Cannot decrypt signed URL", err)
		return err
	}

	debounce := time.Now().Add(time.Second * 30).Unix()

	reportProgress := func(bytes int64) {
		if time.Now().Unix() > debounce {
			fmt.Println("Download bytes", bytes)
			qm.IncreaseLease()
			debounce = time.Now().Add(time.Second * 30).Unix()
		}
	}

	err = azure.DownloadSignedURLToFile(signedURL, body.FileName, reportProgress)
	if err != nil {
		fmt.Println("Failed to download", err)
	}
	fmt.Println("Downloaded ", body.FileName)
	return nil
}
