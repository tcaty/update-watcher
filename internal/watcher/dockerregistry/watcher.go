package dockerregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type image struct {
	name      string
	allowTags string
}

type Watcher struct {
	enabled bool
	name    string
	baseUrl string
	images  []image
}

func NewWatcher(cfg config.Dockerregistry) *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	images := utils.MapArr(cfg.Images, func(i config.Image) image {
		return image{
			name:      i.Name,
			allowTags: i.AllowTags,
		}
	})
	return &Watcher{
		enabled: cfg.Enabled,
		name:    cfg.Name,
		baseUrl: baseUrl,
		images:  images,
	}
}

func (w *Watcher) Enabled() bool {
	return w.enabled
}

func (w *Watcher) Name() string {
	return w.name
}

func (w *Watcher) Targets() []string {
	targets := utils.MapArr(w.images, func(i image) string { return i.name })
	return targets
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

func (w *Watcher) GetLatestVersion(data []byte, target string) (string, error) {
	var tags Tags

	if err := json.Unmarshal(data, &tags); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	for _, t := range tags.Results {
		tag := t.Name
		allowTags := w.getAllowedTagsByName(target)
		match, err := regexp.MatchString(allowTags, tag)
		if err != nil {
			return "", fmt.Errorf("wrong regexp pattern: %v", err)
		}
		if match {
			return tag, nil
		}
	}

	// if there are no tags except latest, only then return it
	return "latest", nil
}

func (w *Watcher) getAllowedTagsByName(name string) string {
	for _, i := range w.images {
		if i.name == name {
			return i.allowTags
		}
	}
	// this case is not possible in general
	// therefore error is useless here
	return ".+"
}
