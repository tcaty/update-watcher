package config

type Config struct {
	Watchers struct {
		Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
		Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
	} `yaml:"watchers"`
}

type Grafanadasboards struct {
	Enabled    bool     `yaml:"enabled"`
	Dashboards []string `yaml:"dashboards"`
}

type Dockerregistry struct {
	Enabled bool `yaml:"enabled"`
	Auth    struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
	Images []string `yaml:"images"`
}
