package config

type Config struct {
	Watchers struct {
		Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
		Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
	} `yaml:"watchers"`
}

type Grafanadasboards struct {
	Enabled    bool     `yaml:"enabled"`
	Name       string   `yaml:"name"`
	Dashboards []string `yaml:"dashboards"`
}

type Dockerregistry struct {
	Enabled bool     `yaml:"enabled"`
	Name    string   `yaml:"name"`
	Images  []string `yaml:"images"`
}
