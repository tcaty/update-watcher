package watcher

import (
	"fmt"
	"log"
	"os"
)

// map target identifire to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type Versions = map[string]string
type Watcher interface {
	IsEnabled() bool
	GetName() string
	Initialize() error
	GetLatestVersions() (Versions, error)
}

func Start(w Watcher) {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", w.GetName()), log.Ldate|log.Ltime)

	logger.Println("Reading configuration...")
	if !w.IsEnabled() {
		logger.Println("Watcher is disabled.")
		return
	}
	logger.Println("Watcher is enabled.")

	if err := w.Initialize(); err != nil {
		logger.Printf("error occured while initialzing watcher: %v\n", err)
		return
	}
	logger.Println("Watcher has been initialized successfully!")
	fmt.Println(w.GetLatestVersions())
	fmt.Println()
}
