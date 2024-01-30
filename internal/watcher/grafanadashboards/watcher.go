package grafanadashboards

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/internal/entities"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Watcher struct {
	enabled    bool
	name       string
	baseUrl    string
	dashboards []dashboard
	embed      *config.Embed
}

type dashboard struct {
	name string
	id   string
}

func NewWatcher(cfg config.Grafanadasboards) *Watcher {
	dashboards := utils.MapArr(cfg.Dashboards, func(d config.Dashboard) dashboard {
		return dashboard{
			name: d.Name,
			id:   d.Id,
		}
	})
	return &Watcher{
		enabled:    cfg.Enabled,
		name:       cfg.Name,
		baseUrl:    "https://grafana.com",
		dashboards: dashboards,
		embed:      &cfg.Embed,
	}
}

func (wt *Watcher) Enabled() bool {
	return wt.enabled
}

func (wt *Watcher) FetchLatestVersionRecords() ([]entities.VersionRecord, error) {
	vrs := make([]entities.VersionRecord, 0, len(wt.targets()))

	for _, target := range wt.targets() {
		url := wt.createUrl(target)

		var revisions Revisions
		var _, err = req.C().R().
			SetSuccessResult(&revisions).
			Get(url)

		if err != nil {
			return nil, err
		}

		vr := entities.VersionRecord{
			Target:  target,
			Version: getLatestRevision(revisions),
		}
		vrs = append(vrs, vr)
	}

	return vrs, nil
}

func (wt *Watcher) CreateMessageAboutUpdates(vrs []entities.VersionRecord) core.Message {
	hrefs := createHrefs(vrs, wt.dashboards)
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
	targets := utils.MapArr(wt.dashboards, func(d dashboard) string { return d.id })
	return targets
}

func (wt *Watcher) createUrl(dashboard string) string {
	return fmt.Sprintf("%s/api/dashboards/%s/revisions", wt.baseUrl, dashboard)
}
