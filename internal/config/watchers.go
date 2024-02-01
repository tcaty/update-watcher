package config

import "github.com/spf13/viper"

// -- Watchers configurations

type Watchers struct {
	Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
	Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
}

type Embed struct {
	Avatar string `yaml:"avatar"`
	Color  int    `yaml:"color"`
	Text   string `yaml:"text"`
}

// -- Grafanadashboards watcher configuration

type Grafanadasboards struct {
	Enabled    bool        `yaml:"enabled"`
	Name       string      `yaml:"name"`
	Dashboards []Dashboard `yaml:"dashboards"`
	Embed      Embed       `yaml:"embed"`
}

type Dashboard struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

func setDetaultGrafanadashboardsWatcherValues() {
	viper.SetDefault("watchers.grafanadasboards.enabled", false)
	viper.SetDefault("watchers.grafanadasboards.name", "grafanadasboards")
	viper.SetDefault("watchers.grafanadasboards.embed.avatar", "https://cdn.icon-icons.com/icons2/2699/PNG/512/grafana_logo_icon_171048.png")
	viper.SetDefault("watchers.grafanadasboards.embed.color", 16296468)
	viper.SetDefault("watchers.grafanadasboards.embed.text", "New revesions released! Checkout:")
}

// -- Dockerregistry watcher configuration

type Dockerregistry struct {
	Enabled bool    `yaml:"enabled"`
	Name    string  `yaml:"name"`
	Images  []Image `yaml:"images"`
	Embed   Embed   `yaml:"embed"`
}

type Image struct {
	Name      string `yaml:"name"`
	AllowTags string `yaml:"allowTags"`
}

func setDefaultDockerregistryWatcherValues() {
	viper.SetDefault("watchers.dockerregistry.name", "dockerregistry")
	viper.SetDefault("watchers.dockerregistry.enabled", false)
	viper.SetDefault("watchers.dockerregistry.images[].allowTags", ".+")
	viper.SetDefault("watchers.dockerregistry.embed.avatar", "https://cdn4.iconfinder.com/data/icons/logos-and-brands/512/97_Docker_logo_logos-512.png")
	viper.SetDefault("watchers.dockerregistry.embed.color", 242424)
	viper.SetDefault("watchers.dockerregistry.embed.text", "New tags released! Checkout:")
}
