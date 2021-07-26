package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the request command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Azure MFT service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
