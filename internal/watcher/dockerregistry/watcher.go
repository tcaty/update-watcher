package dockerregistry

import (
	"fmt"
	"regexp"

	"github.com/imroc/req/v3"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/internal/entities"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Watcher struct {
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
		enabled: cfg.Enabled,
		name:    cfg.Name,
		baseUrl: baseUrl,
		images:  images,
		embed:   &cfg.Embed,
	}
}

func (wt *Watcher) Enabled() bool {
	return wt.enabled
}

func (wt *Watcher) FetchLatestVersionRecords() ([]entities.VersionRecord, error) {
	vrs := make([]entities.VersionRecord, 0, len(wt.targets()))

	for _, target := range wt.targets() {
		url, err := wt.createUrl(target)
		if err != nil {
			return nil, err
		}

		var tags Tags
		_, err = req.C().R().
			SetSuccessResult(&tags).
			Get(url)

		if err != nil {
			return nil, err
		}

		for _, r := range tags.Results {
			tag := r.Name
			allowedTags := wt.getAllowedTagsByName(tag)
			match, err := regexp.MatchString(allowedTags, tag)
			if err != nil {
				return nil, err
			}

			if match {
				vr := entities.VersionRecord{
					Target:  target,
					Version: tag,
				}
				vrs = append(vrs, vr)
				break
			}
		}
	}

	return vrs, nil
}

func (wt *Watcher) CreateMessageAboutUpdates(vrs []entities.VersionRecord) core.Message {
	hrefs := createHrefs(vrs)
	ul := markdown.CreateUnorderedList(hrefs)
	descr := fmt.Sprintf("%s\n%s", wt.embed.Text, ul)
	msg := core.Message{
		Author:      wt.name,
		Avatar:      wt.embed.Avatar,
		Description: descr,
		Color:       wt.embed.Color,
	}
	return msg
}

func (wt *Watcher) targets() []string {
	ts := utils.MapArr(wt.images, func(i image) string { return i.name })
	return ts
}

func (wt *Watcher) createUrl(image string) (string, error) {
	ns, repo, err := splitNsAndRepo(image)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags", wt.baseUrl, ns, repo)
	return url, nil
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
