package daemon

import (
	"encoding/json"
	"fmt"

	"github.com/willhackett/azure-mft/pkg/constant"
)


func handleFileRequest(m constant.Message) error {
	body := constant.FileRequestMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}

	fmt.Println("File Request", body)
	return nil
}

func handleFileHandshake(m constant.Message) error {
	body := constant.FileHandshakeMessage{}
	if err := json.Unmarshal(m.Payload, &body); err != nil {
		return err
	}
	fmt.Println("File Handshake", body)
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