package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     Logger     `yaml:"logger"`
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

	err := viper.Unmarshal(&cfg)
	return cfg, err
}

func setDefaultValues() {
	// -- crontab
	viper.SetDefault("crontab.crontab", "0 */12 * * *")
	viper.SetDefault("crontab.withSeconds", false)
	viper.SetDefault("crontab.execImmediate", false)

	// -- watchers
	// grafanadashboards
	viper.SetDefault("watchers.grafanadasboards.name", "grafanadasboards")
	viper.SetDefault("watchers.grafanadasboards.enabled", false)
	viper.SetDefault(
		"watchers.grafanadasboards.embed.avatar",
		"https://cdn.icon-icons.com/icons2/2699/PNG/512/grafana_logo_icon_171048.png",
	)
	viper.SetDefault("watchers.grafanadasboards.embed.color", 16296468)
	viper.SetDefault("watchers.grafanadasboards.embed.text", "New revesions released! Checkout:")
	// dockerregistry
	viper.SetDefault("watchers.dockerregistry.name", "dockerregistry")
	viper.SetDefault("watchers.dockerregistry.enabled", false)
	viper.SetDefault("watchers.dockerregistry.images[].allowTags", ".+")
	viper.SetDefault(
		"watchers.dockerregistry.embed.avatar",
		"https://cdn4.iconfinder.com/data/icons/logos-and-brands/512/97_Docker_logo_logos-512.png",
	)
	viper.SetDefault("watchers.dockerregistry.embed.color", 242424)
	viper.SetDefault("watchers.dockerregistry.embed.text", "New tags released! Checkout:")

	// -- webhooks
	// discord
	viper.SetDefault("webhooks.discord.enabled", false)
	viper.SetDefault("webhooks.discord.name", "discord")
	// slack
	viper.SetDefault("webhooks.slack.enabled", false)
	viper.SetDefault("webhooks.slack.name", "slack")

	// -- logger
	viper.SetDefault("logger.logLevel", slog.LevelInfo.String())
}
