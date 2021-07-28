package messaging

import (
	"encoding/json"

	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
)

func SendMessage(id string, messageType string, payload []byte, destinationAgent string) error {
	var uuid string
	var err error
	var body []byte

	message := &constant.Message{
		ID:      uuid,
		KeyID:   config.GetKeys().KeyID,
		Agent:   config.GetConfig().Agent.Name,
		Type:    messageType,
		Payload: payload,
	}

	if err := keys.SignMessage(message); err != nil {
		return err
	}

	if body, err = json.Marshal(message); err != nil {
		return err
	}

	return azure.PostMessage(destinationAgent, string(body))
}
