package watcher

import "fmt"

// map target identifire to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type Versions = map[string]string
type Watcher interface {
	IsEnabled() bool
	Initialize() error
	GetLatestVersions() (Versions, error)
}

func Start(w Watcher) {
	if !w.IsEnabled() {
		return
	}
	if err := w.Initialize(); err != nil {
		fmt.Printf("error occured while initialzing watcher: %v", err)
		return
	}
	fmt.Println(w.GetLatestVersions())
}
