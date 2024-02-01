package config

import "github.com/spf13/viper"

// -- Postgres configuration

type Postgres struct {
	Database string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func bindPostgresEnv() {
	viper.BindEnv("postgres.db", "POSTGRES_DB")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
}

func setDefaultPostgresValues() {
	viper.SetDefault("postgres.db", "update-watcher")
	viper.SetDefault("postgres.user", "update-watcher")
	viper.SetDefault("postgres.password", "changeme")
	viper.SetDefault("postgres.host", "0.0.0.0")
	viper.SetDefault("postgres.port", 5432)
}
