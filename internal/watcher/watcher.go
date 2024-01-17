package watcher

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tcaty/update-watcher/pkg/markdown"
)

// map target to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type VersionRecords = map[string]string
type Watcher interface {
	Slog() *slog.Logger
	Enabled() bool
	Name() string
	Targets() []string
	CreateUrl(target string) (string, error)
	CreateHref(target string, version string) *markdown.Href
	GetLatestVersion(data []byte, target string) (string, error)
}

func GetLatestVersions(w Watcher, targets []string) (VersionRecords, error) {
	versionRecords := make(VersionRecords, len(targets))

	for _, t := range targets {
		url, err := w.CreateUrl(t)
		if err != nil {
			return nil, fmt.Errorf("cannot create url: %v", err)
		}

		r, err := getLatestVersion(w, url, t)
		if err != nil {
			return nil, err
		}

		versionRecords[t] = r
	}

	return versionRecords, nil
}

func getLatestVersion(w Watcher, url string, target string) (string, error) {
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
