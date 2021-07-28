package tasks

import (
	"encoding/json"

	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/messaging"
)

func SendFileRequest(sourceFileName string, sourceAgent string, destinationAgent string, destinationFileName string) error {
	var payload []byte
	var err error
	var uuid string

	if uuid, err = constant.GetUUID(); err != nil {
		return err
	}

	details := constant.FileRequestMessage{
		FileName:            sourceFileName,
		DestinationAgent:    destinationAgent,
		DestinationFileName: destinationFileName,
	}

	if payload, err = json.Marshal(details); err != nil {
		return err
	}

	if err = messaging.SendMessage(uuid, constant.FileRequestMessageType, payload, sourceAgent); err != nil {
		return err
	}

	return nil
}

func SendFileHandshake(id string, fileName string, fileSize int64, destinationAgent string) error {
	var payload []byte
	var err error

	if payload, err = json.Marshal(constant.FileHandshakeMessage{
		FileName: fileName,
		FileSize: fileSize,
	}); err != nil {
		return err
	}

	if err = messaging.SendMessage(id, constant.FileHandshakeMessageType, payload, destinationAgent); err != nil {
		return err
	}

	return nil
}

func SendFileHandshakeResponse(id string, accepted bool, destinationAgent string, reason string) error {
	var payload []byte
	var err error

	if payload, err = json.Marshal(constant.FileHandshakeResponseMessage{
		Accepted: accepted,
		Reason:   reason,
	}); err != nil {
		return err
	}

	if err = messaging.SendMessage(id, constant.FileHandshakeResponseMessageType, payload, destinationAgent); err != nil {
		return err
	}
	return nil
}

func SendFileAvailable(id string, signedURL string, fileName string, destinationAgent string) error {
	var payload []byte
	var err error

	if payload, err = json.Marshal(constant.FileAvailableMessage{
		SignedURL: signedURL,
		FileName:  fileName,
	}); err != nil {
		return err
	}

	if err = messaging.SendMessage(id, constant.FileAvailableMessageType, payload, destinationAgent); err != nil {
		return err
	}
	return nil
}
