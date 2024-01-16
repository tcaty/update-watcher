package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	CronJob    CronJob    `yaml:"cronjob"`
	Watchers   Watchers   `yaml:"watchers"`
	Postgresql Postgresql `yaml:"postgresql"`
	Webhooks   Webhooks   `yaml:"webhooks"`
}

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
	setDefaultValues()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error occured during reading config file: %v", err)
	}
	slog.Info("Reading config...", "configFile", viper.ConfigFileUsed())

	err := viper.Unmarshal(&cfg)
	return cfg, err
}

func setDefaultValues() {
	// -- crontab
	viper.SetDefault("crontab.crontab", "0 */12 * * *")
	viper.SetDefault("crontab.withSeconds", false)
	// -- watchers
	viper.SetDefault("watchers.grafanadasboards.name", "grafanadasboards")
	viper.SetDefault("watchers.dockerregistry.name", "dockerregistry")
	viper.SetDefault("watchers.dockerregistry.enabled", false)
	viper.SetDefault("watchers.grafanadasboards.name", false)
	viper.SetDefault("watchers.dockerregistry.images[].allowTags", ".+")
	// -- webhooks
	viper.SetDefault("webhooks.discord.enabled", false)
	viper.SetDefault("webhooks.slack.enabled", false)
}
