package config

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/willhackett/azure-mft/pkg/keys"
)

type AgentConf struct {
	Name string `mapstructure:"name"`
}

type PathsConf struct {
	KeysDir string	`mapstructure:"keys_dir"`
	CacheDir string `mapstructure:"cache_dir"`
	TmpDir	string `mapstructure:"tmp_dir"`
}

type AzureConf struct {
	AccountName string	`mapstructure:"account_name"`
	AccountKey string	`mapstructure:"account_key"`
}

type Exit struct {
	AgentName string `mapstructure:"agent_name"`
	FileMatch string `mapstructure:"file_match"`
	Command	 string	`mapstructure:"command"`
}

type Config struct {
	Agent AgentConf `mapstructure:"agent"`

	Paths PathsConf `mapstructure:"paths"`
	
	Azure AzureConf `mapstructure:"azure"`
	
	Exits []Exit `mapstructure:"exits"`
}

var (
	ConfigFilePath	string

	config	Config

	privateKey	*rsa.PrivateKey

	publicKey	crypto.PublicKey
)

func Init() {
	if ConfigFilePath != "" {
		viper.SetConfigFile(ConfigFilePath)
	} else {
		userConfigDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(userConfigDir)
    viper.AddConfigPath("/var/azmft/")
    viper.AddConfigPath("./.var/azmft/")
    
		viper.SetConfigType("yaml")
		viper.SetConfigName("azmft.config.yaml")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Config loaded: ", viper.ConfigFileUsed())
	}

	if err := viper.UnmarshalKey("config", &config); err != nil {
		cobra.CheckErr(err)
	}

	if config.Agent.Name == "" {
		cobra.CheckErr(errors.New("config.agent.name is not specified"))
	}
	if config.Agent.Name == "publickeys" {
		cobra.CheckErr(errors.New("publickeys is a reserved agent name"))
	}
	if config.Azure.AccountName == "" {
		cobra.CheckErr(errors.New("config.azure.account_name is not specified"))
	}
	if config.Azure.AccountKey == "" {
		cobra.CheckErr(errors.New("config.azure.account_key is not specified"))
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}

	if config.Paths.CacheDir == "" {
		config.Paths.CacheDir = cacheDir + "/azmft"
	}
	if config.Paths.KeysDir == "" {
		config.Paths.KeysDir = userHomeDir + "/azmft/keys"
	}
	if config.Paths.TmpDir == "" {
		config.Paths.TmpDir = os.TempDir()
	}

	privateKey, publicKey, err = keys.GetKeys(config.Paths.KeysDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	fmt.Print(config, publicKey)
}

func GetConfig() Config {
	return config
}

func GetKeys() (*rsa.PrivateKey, interface{}) {
	return privateKey, publicKey
}
