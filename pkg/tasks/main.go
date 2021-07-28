package tasks

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/logger"
	"github.com/willhackett/azure-mft/pkg/messaging"
)

func SendFileRequest(sourceFileName string, sourceAgent string, destinationAgent string, destinationFileName string) error {
	var payload []byte
	var err error
	uuid, _ := constant.GetUUID()
	log := logger.Get().WithFields(logrus.Fields{
		"id":                  uuid,
		"event":               "SendFileRequest",
		"sourceFileName":      sourceFileName,
		"sourceAgent":         sourceAgent,
		"destinationAgent":    destinationAgent,
		"destinationFileName": destinationFileName,
	})

	details := constant.FileRequestMessage{
		FileName:            sourceFileName,
		DestinationAgent:    destinationAgent,
		DestinationFileName: destinationFileName,
	}

	if payload, err = json.Marshal(details); err != nil {
		log.Trace(err)
		return err
	}

	if err = messaging.SendMessage(uuid, constant.FileRequestMessageType, payload, sourceAgent); err != nil {
		log.Trace(err)
		return err
	}

	log.Info("Successfully sent file request")

	return nil
}

func SendFileHandshake(id string, fileName string, fileSize int64, destinationAgent string) error {
	var payload []byte
	var err error
	log := logger.Get().WithFields(logrus.Fields{
		"id":               id,
		"event":            "SendFileHandshake",
		"destinationAgent": destinationAgent,
		"fileName":         fileName,
		"fileSize":         fileSize,
	})

	if payload, err = json.Marshal(constant.FileHandshakeMessage{
		FileName: fileName,
		FileSize: fileSize,
	}); err != nil {
		log.Trace(err)
		return err
	}

	if err = messaging.SendMessage(id, constant.FileHandshakeMessageType, payload, destinationAgent); err != nil {
		log.Error("Failed to send file handshake", err)
		log.Trace(err)
		return err
	}

	log.Info("Successfully sent file handshake")

	return nil
}

func SendFileHandshakeResponse(id string, accepted bool, destinationAgent string, reason string) error {
	var payload []byte
	var err error
	log := logger.Get().WithFields(logrus.Fields{
		"id":               id,
		"event":            "SendFileHandshakeResponse",
		"accepted":         accepted,
		"destinationAgent": destinationAgent,
	})

	if payload, err = json.Marshal(constant.FileHandshakeResponseMessage{
		Accepted: accepted,
		Reason:   reason,
	}); err != nil {
		log.Trace(err)
		return err
	}

	if err = messaging.SendMessage(id, constant.FileHandshakeResponseMessageType, payload, destinationAgent); err != nil {
		log.Error("Failed to send file handshake response", err)
		return err
	}

	log.Info("Successfully sent file handshake response")

	return nil
}

func SendFileAvailable(id string, signedURL string, fileName string, destinationAgent string) error {
	var payload []byte
	var err error
	log := logger.Get().WithFields(logrus.Fields{
		"id":               id,
		"event":            "SendFileAvailable",
		"destinationAgent": destinationAgent,
		"fileName":         fileName,
	})

	if payload, err = json.Marshal(constant.FileAvailableMessage{
		SignedURL: signedURL,
		FileName:  fileName,
	}); err != nil {
		log.Trace(err)
		return err
	}

	if err = messaging.SendMessage(id, constant.FileAvailableMessageType, payload, destinationAgent); err != nil {
		log.Error("Failed to send file available", err)
		return err
	}

	log.Info("Successfully sent file available")

	return nil
}
