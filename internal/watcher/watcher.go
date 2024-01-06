package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tcaty/update-watcher/internal/repository"
)

// map target to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type VersionRecords = map[string]string
type Watcher[T comparable] interface {
	IsEnabled() bool
	GetName() string
	CreateUrl(target string) (string, error)
	GetTargets() []string
	GetLatestVersion(versions T) string
}

func Initialize[T comparable](w Watcher[T]) {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", w.GetName()), log.Ldate|log.Ltime)

	logger.Println("Reading configuration...")
	if !w.IsEnabled() {
		logger.Println("Watcher is disabled.")
		return
	}

	// TODO: implement auth logic here

	logger.Println("Watcher is enabled.")
	logger.Println("Watcher has been initialized successfully!")
}

func Tick[T comparable](w Watcher[T], r *repository.Repository) {
	targets := w.GetTargets()
	versions, err := getLatestVersions[T](w, targets)
	if err != nil {
		return
	}
	for target, version := range versions {
		updated, err := r.UpdateVersionRecord(target, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v", err)
		}
		if updated {
			fmt.Printf("updated %s -> %s\n", target, version)
		}
	}
}

func getLatestVersions[T comparable](w Watcher[T], targets []string) (VersionRecords, error) {
	versionRecords := make(VersionRecords, len(targets))

	for _, t := range targets {
		url, err := w.CreateUrl(t)
		if err != nil {
			return nil, fmt.Errorf("cannot create url: %v", err)
		}

		r, err := getLatestVersion[T](w, url)
		if err != nil {
			return nil, err
		}

		versionRecords[t] = r
	}

	return versionRecords, nil
}

func getLatestVersion[T comparable](w Watcher[T], url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %v", err)
	}

	var versions T
	if err := json.Unmarshal(body, &versions); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	latestVersion := w.GetLatestVersion(versions)
	return latestVersion, nil
}
