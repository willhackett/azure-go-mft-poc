package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
	"github.com/willhackett/azure-mft/pkg/keys"
)

// copyCmd represents the copy command
var (
	destinationAgent string
	destinationFileName  string
	fileName				 string

	copyCmd = &cobra.Command{
		Use:   "copy",
		Short: "Copy a file to another agent",
		Run: func(cmd *cobra.Command, args []string) {
			if (fileName == "" || destinationFileName == "" || destinationAgent == "") {
				err := errors.New("file name, destination file name and destination agent must be specified")
				cobra.CheckErr(err)
			}

			uuid, err := constant.GetUUID()
			cobra.CheckErr(err)

			payload := constant.FileRequestMessage{
				FileName:			 fileName,
				DestinationAgent: destinationAgent,
				DestinationFileName:  destinationFileName,
			}

			marshalledPayload, err := json.Marshal(payload)
			cobra.CheckErr(err)

			message := &constant.Message{
				ID: uuid,
				KeyID: config.GetKeys().KeyID,
				Type: constant.FileRequestMessageType,
				Agent: config.GetConfig().Agent.Name,
				Payload: marshalledPayload,
			}

			keys.SignMessage(message)

			marshalledMessage, err := json.Marshal(message)
			cobra.CheckErr(err)

			azure.PostMessage(message.Agent, string(marshalledMessage))
			fmt.Println("Marhsalled ", string(marshalledMessage))
		},
	}
)

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.PersistentFlags().StringVar(&destinationAgent, "destinationAgent", "", "Destination agent")
	copyCmd.PersistentFlags().StringVar(&destinationFileName, "destinationFileName", "", "Destination file name")
	copyCmd.PersistentFlags().StringVar(&fileName, "fileName", "", "File name")
}
