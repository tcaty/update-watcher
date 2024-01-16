package config

type Watchers struct {
	Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
	Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
}

type Dashboard struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type Grafanadasboards struct {
	Enabled    bool        `yaml:"enabled"`
	Name       string      `yaml:"name"`
	Dashboards []Dashboard `yaml:"dashboards"`
}

type Image struct {
	Name      string `yaml:"name"`
	AllowTags string `yaml:"allowTags"`
}

type Dockerregistry struct {
	Enabled bool    `yaml:"enabled"`
	Name    string  `yaml:"name"`
	Images  []Image `yaml:"images"`
}
