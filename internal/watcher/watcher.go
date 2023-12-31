package watcher

import "fmt"

// map target identifire to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type Versions = map[string]string
type Watcher interface {
	GetLatestVersions() (Versions, error)
}

func Tick(w Watcher) {
	fmt.Println(w.GetLatestVersions())
}
