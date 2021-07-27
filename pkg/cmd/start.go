package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/daemon"
)

// startCmd represents the request command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Azure MFT service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting agent")
		daemon.Init()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
