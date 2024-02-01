package config

import "github.com/spf13/viper"

// -- Webhooks configuration

type Webhooks struct {
	Discord Discord `yaml:"discord"`
}

// -- Discord webhook configuration

type Discord struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
}

func bindDiscordWebhookEnv() {
	viper.BindEnv("webhooks.discord.url", "WEBHOOKS_DISCORD_URL")
}

func setDefaultDiscordWebhookValues() {
	viper.SetDefault("webhooks.discord.enabled", true)
	viper.SetDefault("webhooks.discord.name", "discord")
	viper.SetDefault("webhooks.discord.url", "")
}
