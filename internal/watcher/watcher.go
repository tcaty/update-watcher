package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// map target identifire to it's latest version
// grafanadashboards: {"1860": "31"}
// dockerregistry: {"grafana/dashboard": "10.7.4"}
type Versions = map[string]string
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
		logger.Println()
		return
	}

	// TODO:

	logger.Println("Watcher is enabled.")
	logger.Println("Watcher has been initialized successfully!")
}

func Tick[T comparable](w Watcher[T]) error {
	targets := w.GetTargets()
	versions, err := getLatestVersions[T](w, targets)
	if err != nil {
		return err
	}
	fmt.Println(versions)
	return nil
}

func getLatestVersions[T comparable](w Watcher[T], targets []string) (Versions, error) {
	versions := make(Versions, len(targets))

	for _, t := range targets {
		url, err := w.CreateUrl(t)
		if err != nil {
			return nil, fmt.Errorf("cannot create url: %v", err)
		}

		r, err := getLatestVersion[T](w, url)
		if err != nil {
			return nil, err
		}

		versions[t] = r
	}

	return versions, nil
}

func getLatestVersion[T comparable](w Watcher[T], url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get grafanadashboards url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read revisions body: %v", err)
	}

	var versions T
	if err := json.Unmarshal(body, &versions); err != nil {
		return "", fmt.Errorf("cannot unmarshal revisions: %v", err)
	}

	latestVersion := w.GetLatestVersion(versions)
	return latestVersion, nil
}
