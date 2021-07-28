package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
	"github.com/willhackett/azure-mft/pkg/logger"
	"github.com/willhackett/azure-mft/pkg/registry"
	"github.com/willhackett/azure-mft/pkg/tasks"
)

func handleFileRequest(m constant.Message) error {
	log := logger.Get().WithFields(logrus.Fields{
		"id":    m.ID,
		"event": "HandleFileRequest",
	})
	body := constant.FileRequestMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

	registry.AddTransfer(m.ID, body, 5*60*60*1000)

	file, err := os.Open(body.FileName)
	if err != nil {
		log.Error("Cannot open file for reading", err)
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error("Cannot read file information", err)
		return err
	}

	fileSize := fileInfo.Size()
	log.Debug(fmt.Sprintf("File size: %d", fileSize))

	return tasks.SendFileHandshake(m.ID, body.DestinationFileName, fileSize, body.DestinationAgent)
}

func handleFileHandshake(m constant.Message) error {
	log := logger.Get().WithFields(logrus.Fields{
		"id":    m.ID,
		"event": "HandleFileHandshake",
	})
	body := constant.FileHandshakeMessage{}
	log.Info("Received file handshake request from " + m.Agent)
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		log.Error("File handshake has invalid payload", err)
		return err
	}

	_, err := os.OpenFile(body.FileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error(fmt.Sprintf("Cannot open destination path: %s", body.FileName), err)
		tasks.SendFileHandshakeResponse(m.ID, false, m.Agent, fmt.Sprintf("Cannot open destination path: %s", body.FileName))
		return nil
	}

	err = tasks.SendFileHandshakeResponse(m.ID, true, m.Agent, "")

	return err
}

func handleFileHandshakeResponse(qm *QueueMessage, m constant.Message) error {
	log := logger.Get().WithFields(logrus.Fields{
		"id":    m.ID,
		"event": "HandleFileHandshakeResponse",
	})
	log.Info(fmt.Sprintf("Received file handshake response from %s", m.Agent))

	body := constant.FileHandshakeResponseMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		log.Error("File handshake response has invalid payload", err)
		return err
	}

	if !body.Accepted {
		log.Warn("File handshake was not accepted, end of transaction.")
		return nil
	}

	transfer, ok := registry.GetTransfer(m.ID)
	if !ok {
		log.Debug("Cannot find transfer", m.ID)
		return errors.New("transfer expired or did not originate from this node")
	}

	debounce := time.Now().Add(time.Second * 30).Unix()

	reportProgress := func(bytes int64) {
		if time.Now().Unix() > debounce {
			qm.IncreaseLease()
			log.Debug(fmt.Sprintf("Uploaded bytes: %d and increased message visibility timeout", bytes))
			debounce = time.Now().Add(time.Second * 30).Unix()
		}
	}

	signedURL, err := azure.UploadFromFile(transfer.Details.DestinationAgent, m.ID, transfer.Details.FileName, reportProgress)
	if err != nil {
		log.Error("Failed to upload file", err)
		return err
	}

	encryptedSignedURL, err := keys.EncryptString(transfer.Details.DestinationAgent, m.KeyID, signedURL)
	if err != nil {
		log.Error("Failed to encrypt signed URL", err)
		return err
	}

	err = tasks.SendFileAvailable(m.ID, encryptedSignedURL, transfer.Details.DestinationFileName, transfer.Details.DestinationAgent)

	return err
}

func handleFileAvailable(qm *QueueMessage, m constant.Message) error {
	log := logger.Get().WithFields(logrus.Fields{
		"id":    m.ID,
		"event": "HandleFileAvailable",
	})
	log.Info(fmt.Sprintf("Received file available from %s", m.Agent))

	body := constant.FileAvailableMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

	signedURL, err := keys.DecryptString(body.SignedURL)
	if err != nil {
		log.Error("Failed to decrypt signed URL", err)
		return err
	}

	debounce := time.Now().Add(time.Second * 30).Unix()

	reportProgress := func(bytes int64) {
		if time.Now().Unix() > debounce {
			qm.IncreaseLease()
			log.Debug(fmt.Sprintf("Downloaded bytes: %d and increased message visibility timeout", bytes))
			debounce = time.Now().Add(time.Second * 30).Unix()
		}
	}

	err = azure.DownloadSignedURLToFile(signedURL, body.FileName, reportProgress)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to download file: %s", body.FileName), err)
	}
	log.Info(fmt.Sprintf("Downloaded file: %s", body.FileName))
	return nil
}
