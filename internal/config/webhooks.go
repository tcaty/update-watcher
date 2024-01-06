package config

type Webhooks struct {
	Discord Discord `yaml:"discord"`
	Slack   Slack   `yaml:"slack"`
}

type Discord struct {
	Enabled bool   `yaml:"enabled"`
	Url     string `yaml:"url"`
	Avatar  string `yaml:"avatar"`
	Author  string `yaml:"author"`
}

type Slack struct {
	Enabled bool `yaml:"enabled"`
}
