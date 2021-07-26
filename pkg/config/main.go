package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Paths struct {
	KeysDir	string
	CacheDir string
	TmpDir	string
}

type AzureResources struct {
	AccountName string
	AccountKey string
	StorageAccountName string
}

var (
	ConfigFilePath	string
)

func Init() {
	if ConfigFilePath != "" {
		viper.SetConfigFile(ConfigFilePath)
	} else {
		userConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(userConfigDir)
    viper.AddConfigPath("/var/azmft/")
    
		viper.SetConfigType("yaml")
		viper.SetConfigName("azmft.config.yaml")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Configuration loaded from ", viper.ConfigFileUsed())
	}
}