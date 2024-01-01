package dockerregistry

import (
	"fmt"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
)

type Watcher struct {
	enabled   bool
	baseUrl   string
	authToken string
}

func NewWatcher(cfg config.Dockerregistry) *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	return &Watcher{
		enabled:   cfg.Enabled,
		baseUrl:   baseUrl,
		authToken: "",
	}
}

func (w *Watcher) IsEnabled() bool {
	return w.enabled
}

func (w *Watcher) Initialize() error {
	// TODO: auth to dockerhub ang get token here
	fmt.Println("Dockerregistry wathcher initialized.")
	return nil
}

func (w *Watcher) GetLatestVersions() (watcher.Versions, error) {
	w.getLatestTag()
	return nil, nil
}

func (w *Watcher) getLatestTag() (string, error) {
	return "", nil
}
