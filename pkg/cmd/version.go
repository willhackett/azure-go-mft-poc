package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Display the version number",
  Long:  `Display the version number`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Version 1.0.0")
  },
}