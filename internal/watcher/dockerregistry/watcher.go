package dockerregistry

import (
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
)

type Watcher struct {
	baseUrl   string
	authToken string
}

func NewWatcher(cfg config.Dockerregistry) *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	// TODO: auth to dockerhub ang get token here
	return &Watcher{
		baseUrl:   baseUrl,
		authToken: "",
	}
}

func (w *Watcher) GetLatestVersions() (watcher.Versions, error) {
	w.getLastTag()
	return nil, nil
}

func (w *Watcher) getLastTag() (string, error) {
	return "", nil
}
