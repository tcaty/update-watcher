package dockerregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Watcher struct {
	slog    *slog.Logger
	enabled bool
	name    string
	baseUrl string
	images  []image
	embed   *config.Embed
}

type image struct {
	name      string
	allowTags string
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
		slog:    slog.Default().With("watcher", cfg.Name),
		enabled: cfg.Enabled,
		name:    cfg.Name,
		baseUrl: baseUrl,
		images:  images,
		embed:   &cfg.Embed,
	}
}

func (wt *Watcher) Slog() *slog.Logger {
	return wt.slog
}

func (wt *Watcher) Enabled() bool {
	return wt.enabled
}

func (wt *Watcher) Name() string {
	return wt.name
}

func (wt *Watcher) Targets() []string {
	targets := utils.MapArr(wt.images, func(i image) string { return i.name })
	return targets
}

func (wt *Watcher) Embed() *config.Embed {
	return wt.embed
}

func (wt *Watcher) CreateUrl(image string) (string, error) {
	b := []byte(image)
	i := bytes.IndexByte(b, byte('/'))

	if i < 0 {
		return "", errors.New("docker image should fit the format {namespace}/{repository}")
	}

	ns, repo := string(b[:i]), string(b[i+1:])
	url := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags", wt.baseUrl, ns, repo)

	return url, nil
}

func (wt *Watcher) CreateHref(target string, version string) *markdown.Href {
	text := fmt.Sprintf("%s:%s", target, version)
	link := fmt.Sprintf("https://hub.docker.com/r/%s/tags", target)
	href := markdown.NewHref(text, link)
	return href
}

func (wt *Watcher) GetLatestVersion(data []byte, target string) (string, error) {
	var tags Tags

	if err := json.Unmarshal(data, &tags); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	for _, t := range tags.Results {
		tag := t.Name
		allowTags := wt.getAllowedTagsByName(target)
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

func (wt *Watcher) getAllowedTagsByName(name string) string {
	for _, i := range wt.images {
		if i.name == name {
			return i.allowTags
		}
	}
	// this case is not possible in general
	// therefore error is useless here
	return ".+"
}
