package watcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tcaty/update-watcher/pkg/markdown"
)

// map target to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type VersionRecords = map[string]string
type Watcher interface {
	IsEnabled() bool
	GetName() string
	CreateUrl(target string) (string, error)
	CreateHref(target string, version string) *markdown.Href
	GetTargets() []string
	GetLatestVersion(data []byte) (string, error)
}

func Initialize(w Watcher) error {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", w.GetName()), log.Ldate|log.Ltime)

	logger.Println("Reading configuration...")
	if !w.IsEnabled() {
		logger.Println("Watcher is disabled.")
		return nil
	}

	// TODO: implement auth logic here
	// somewhere here return error

	logger.Println("Watcher is enabled.")
	logger.Println("Watcher has been initialized successfully!")

	return nil
}

func GetLatestVersions(w Watcher, targets []string) (VersionRecords, error) {
	versionRecords := make(VersionRecords, len(targets))

	for _, t := range targets {
		url, err := w.CreateUrl(t)
		if err != nil {
			return nil, fmt.Errorf("cannot create url: %v", err)
		}

		r, err := getLatestVersion(w, url)
		if err != nil {
			return nil, err
		}

		versionRecords[t] = r
	}

	return versionRecords, nil
}

func getLatestVersion(w Watcher, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %v", err)
	}

	latestVersion, err := w.GetLatestVersion(body)
	if err != nil {
		return "", fmt.Errorf("cannot get latest version: %v", err)
	}

	return latestVersion, nil
}
