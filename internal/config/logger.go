package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

// -- Logger configuration

type Logger struct {
	LogLevel string `yaml:"logLevel"`
}

func setDefaultLoggerValues() {
	viper.SetDefault("logger.logLevel", slog.LevelInfo.String())
}
