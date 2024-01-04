package config

type Config struct {
	Watchers   Watchers   `yaml:"watchers"`
	Postgresql Postgresql `yaml:"postgresql"`
}

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

type Postgresql struct {
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
