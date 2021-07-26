package cmd

import (
	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azmft",
	Short: "Azure Managed File Transfer",
	Long: `Managed File Transfer between agents with Azure as a transport.

For more information visit https://github.com/willhackett/azure-mft`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(config.Init)

	rootCmd.PersistentFlags().StringVar(&config.ConfigFilePath, "config", "", "config file location (default is ~/.config/azmft.config.yaml")
}
