package config

type Webhooks struct {
	Discord Discord `yaml:"discord"`
	Slack   Slack   `yaml:"slack"`
}

type Discord struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
}

type Slack struct {
	Enabled bool   `yaml:"enabled"`
	Name    string `yaml:"name"`
}
