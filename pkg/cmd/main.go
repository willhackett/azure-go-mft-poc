package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "azmft",
  Short: "Azure MFT is managed file transfer using just Azure Blob & Queue storage",
  Long: `A fast and simple MFT agent that hinges on Azure - see https://github.com/willhackett/azure-mft for more details`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Stuff")
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}