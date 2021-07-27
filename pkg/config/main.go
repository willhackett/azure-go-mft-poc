package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type AllowFilesFrom []string

type AllowRequestsFrom []string

type Keys struct {
	KeyID string
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type Config struct {
	Agent AgentConf `mapstructure:"agent"`

	Paths PathsConf `mapstructure:"paths"`
	
	Azure AzureConf `mapstructure:"azure"`
	
	Exits []Exit `mapstructure:"exits"`

	AllowFilesFrom AllowFilesFrom `mapstructure:"allow_files_from"`

	AllowRequestsFrom AllowRequestsFrom `mapstructure:"allow_requests_from"`
}

var (
	ConfigFilePath	string

	config	Config

	privateKey	*rsa.PrivateKey

	publicKey *rsa.PublicKey

	keyID		string
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
	cobra.CheckErr(err)
	userHomeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if config.Paths.CacheDir == "" {
		config.Paths.CacheDir = cacheDir + "/azmft"
	}
	if config.Paths.KeysDir == "" {
		config.Paths.KeysDir = userHomeDir + "/azmft/keys"
	}
	if config.Paths.TmpDir == "" {
		config.Paths.TmpDir = os.TempDir()
	}
}

func GetConfig() Config {
	return config
}

func GetKeys() Keys {
	return Keys{
		KeyID: keyID,
		PublicKey: publicKey,
		PrivateKey: privateKey,
	}
}

func SetKeys(privateKeyIn *rsa.PrivateKey, publicKeyIn *rsa.PublicKey, keyIDIn string) {
	publicKey = publicKeyIn
	privateKey = privateKeyIn
	keyID = keyIDIn
}