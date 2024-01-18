package watcher

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

type Watcher interface {
	Slog() *slog.Logger
	Enabled() bool
	Name() string
	Targets() []string
	Embed() *config.Embed
	CreateUrl(target string) (string, error)
	CreateHref(target string, version string) *markdown.Href
	GetLatestVersion(data []byte, target string) (string, error)
}

// map target unique identifier to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type VersionRecords = map[string]string

func FetchLatestVersionRecords(w Watcher, targets []string) (VersionRecords, error) {
	versionRecords := make(VersionRecords, len(targets))

	for _, t := range targets {
		url, err := w.CreateUrl(t)
		if err != nil {
			return nil, fmt.Errorf("cannot create url: %v", err)
		}

		v, err := fetchLatestVersion(w, url, t)
		if err != nil {
			return nil, fmt.Errorf("cannot fetch latest version %s: %v", url, err)
		}

		versionRecords[t] = v
	}

	return versionRecords, nil
}

func fetchLatestVersion(w Watcher, url string, target string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %v", err)
	}

	latestVersion, err := w.GetLatestVersion(body, target)
	if err != nil {
		return "", fmt.Errorf("cannot get latest version: %v", err)
	}

	return latestVersion, nil
}
