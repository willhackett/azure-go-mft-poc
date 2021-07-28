package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/logger"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the current version",
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetApp("Version")
		logger.Get().Info("Version " + config.Version)
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)
}
