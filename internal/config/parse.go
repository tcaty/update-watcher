package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Parse(cfgFile string) (*Config, error) {
	var cfg *Config

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("can't get user home dir: %v", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error occured during reading config file: %v", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	err := viper.Unmarshal(&cfg)
	return cfg, err
}
