package config

type Watchers struct {
	Grafanadasboards Grafanadasboards `yaml:"grafanadasboards"`
	Dockerregistry   Dockerregistry   `yaml:"dockerregistry"`
}

type Grafanadasboards struct {
	Enabled    bool        `yaml:"enabled"`
	Name       string      `yaml:"name"`
	Dashboards []Dashboard `yaml:"dashboards"`
	Embed      Embed       `yaml:"embed"`
}

type Dockerregistry struct {
	Enabled bool    `yaml:"enabled"`
	Name    string  `yaml:"name"`
	Images  []Image `yaml:"images"`
	Embed   Embed   `yaml:"embed"`
}

type Embed struct {
	Avatar string `yaml:"avatar"`
	Color  int    `yaml:"color"`
}

type Dashboard struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type Image struct {
	Name      string `yaml:"name"`
	AllowTags string `yaml:"allowTags"`
}
