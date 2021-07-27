package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/tasks"
)

// copyCmd represents the copy command
var (
	destinationAgent    string
	destinationFileName string
	fileName            string

	copyCmd = &cobra.Command{
		Use:   "copy",
		Short: "Copy a file to another agent",
		Run: func(cmd *cobra.Command, args []string) {
			if fileName == "" || destinationFileName == "" || destinationAgent == "" {
				err := errors.New("file name, destination file name and destination agent must be specified")
				cobra.CheckErr(err)
			}

			workingDir, err := os.Getwd()
			cobra.CheckErr(err)

			if fileName[0:1] != "/" {
				fileName = path.Join(workingDir, fileName)
			}

			if destinationFileName != "/" {
				cobra.CheckErr(errors.New("destination file namemust be an absolute path"))
			}

			if err = tasks.SendFileRequest(fileName, config.GetConfig().Agent.Name, destinationAgent, destinationFileName); err != nil {
				cobra.CheckErr(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.PersistentFlags().StringVar(&destinationAgent, "destinationAgent", "", "Destination agent")
	copyCmd.PersistentFlags().StringVar(&destinationFileName, "destinationFileName", "", "Destination file name")
	copyCmd.PersistentFlags().StringVar(&fileName, "fileName", "", "File name")
}
