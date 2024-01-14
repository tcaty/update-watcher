package dockerregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

type Watcher struct {
	enabled bool
	name    string
	baseUrl string
	images  []string
}

func NewWatcher(cfg config.Dockerregistry) *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	return &Watcher{
		enabled: cfg.Enabled,
		name:    cfg.Name,
		baseUrl: baseUrl,
		images:  cfg.Images,
	}
}

func (w *Watcher) IsEnabled() bool {
	return w.enabled
}

func (w *Watcher) GetName() string {
	return w.name
}

func (w *Watcher) GetTargets() []string {
	return w.images
}

func (w *Watcher) CreateUrl(image string) (string, error) {
	b := []byte(image)
	i := bytes.IndexByte(b, byte('/'))

	if i < 0 {
		return "", errors.New("docker image should fit the format {namespace}/{repository}")
	}

	ns, repo := string(b[:i]), string(b[i+1:])
	url := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags", w.baseUrl, ns, repo)

	return url, nil
}

func (w *Watcher) CreateHref(target string, version string) *markdown.Href {
	text := fmt.Sprintf("%s:%s", target, version)
	link := fmt.Sprintf("https://hub.docker.com/r/%s/tags", target)
	href := markdown.NewHref(text, link)
	return href
}

func (w *Watcher) GetLatestVersion(data []byte) (string, error) {
	var tags Tags

	if err := json.Unmarshal(data, &tags); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	for _, t := range tags.Results {
		name := t.Name
		if name != "latest" {
			return name, nil
		}
	}

	// if there are no tags except latest, only then return it
	return "latest", nil
}
