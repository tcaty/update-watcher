package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   Logger   `yaml:"logger"`
	CronJob  CronJob  `yaml:"cronjob"`
	Watchers Watchers `yaml:"watchers"`
	Postgres Postgres `yaml:"postgres"`
	Webhooks Webhooks `yaml:"webhooks"`
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

	setDefaultValues()
	bindEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error occured during reading config file: %v", err)
	}

	err := viper.Unmarshal(&cfg)
	fmt.Println("prefix", viper.GetEnvPrefix(), os.Getenv("POSTGRES_HOST"))
	fmt.Println(os.Getenv("POSTGRES_USER"))

	return cfg, err
}

func setDefaultValues() {
	setDefaultCronJobValues()
	setDefaultPostgresValues()
	setDefaultLoggerValues()

	setDetaultGrafanadashboardsWatcherValues()
	setDefaultDockerregistryWatcherValues()

	setDefaultDiscordWebhookValues()
}

func bindEnv() {
	bindDiscordWebhookEnv()
	bindPostgresEnv()
}

func init() {
	viper.AutomaticEnv()
}
