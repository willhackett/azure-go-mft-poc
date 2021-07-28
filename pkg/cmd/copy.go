package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/logger"
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
			log := logger.Get()

			if fileName == "" || destinationFileName == "" || destinationAgent == "" {
				log.Fatal("File name, destination file ame and destination agent must be specified")
				os.Exit(1)
			}

			workingDir, err := os.Getwd()
			if err != nil {
				log.Fatal("Cannot determine working directory")
				log.Trace(err)
				os.Exit(1)
			}

			if fileName[0:1] != "/" {
				fileName = path.Join(workingDir, fileName)
			}

			if destinationFileName[0:1] != "/" {
				log.Fatal("The destination filename must have an absolute path")
				os.Exit(1)
			}

			if err = tasks.SendFileRequest(fileName, config.GetConfig().Agent.Name, destinationAgent, destinationFileName); err != nil {
				log.Fatal("Cannot copy file")
				log.Trace(err)
				os.Exit(1)
			}

			log.Info("Done")
		},
	}
)

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.PersistentFlags().StringVar(&destinationAgent, "destinationAgent", "", "Destination agent")
	copyCmd.PersistentFlags().StringVar(&destinationFileName, "destinationFileName", "", "Destination file name")
	copyCmd.PersistentFlags().StringVar(&fileName, "fileName", "", "File name")
}
