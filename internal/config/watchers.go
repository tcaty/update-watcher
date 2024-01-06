package config

type Watchers struct {
	Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
	Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
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
