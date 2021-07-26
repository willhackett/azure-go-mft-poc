package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(requestCmd)
}
