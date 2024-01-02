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

	viper.SetDefault("watchers.grafanadasboards.name", "grafanadasboards")
	viper.SetDefault("watchers.dockerregistry.name", "dockerregistry")
	viper.SetDefault("watchers.dockerregistry.enabled", false)
	viper.SetDefault("watchers.grafanadasboards.name", false)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error occured during reading config file: %v", err)
	}
	fmt.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())

	err := viper.Unmarshal(&cfg)
	return cfg, err
}
