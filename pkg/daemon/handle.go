package daemon

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/tasks"
)

func handleFileRequest(m constant.Message) error {
	body := constant.FileRequestMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

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

	_, err := os.Open(body.FileName)
	if err != nil {
		tasks.SendFileHandshakeResponse(m.ID, false, fmt.Sprintf("Cannot open destination path: %s", err), m.Agent)
		return nil
	}

	if err = tasks.SendFileHandshakeResponse(m.ID, true, "", m.Agent); err != nil {
		return err
	}

	return nil
}

func handleFileHandshakeResponse(m constant.Message) error {
	body := constant.FileHandshakeResponseMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}
	fmt.Println("File Handshake Response", body)
	return nil

}

func handleFileAvailable(m constant.Message) error {
	body := constant.FileAvailableMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}
	fmt.Println("File Available", body)
	return nil

}
