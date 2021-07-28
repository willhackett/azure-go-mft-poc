package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/daemon"
	"github.com/willhackett/azure-mft/pkg/logger"
)

// startCmd represents the request command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Azure MFT service",
	Run: func(cmd *cobra.Command, args []string) {
		logger.SetApp("Daemon")
		logger.Get().Info("Agent started")
		daemon.Init()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
