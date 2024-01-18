package grafanadashboards

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Watcher struct {
	slog       *slog.Logger
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
		slog:       slog.Default().With("watcher", cfg.Name),
		enabled:    cfg.Enabled,
		name:       cfg.Name,
		baseUrl:    "https://grafana.com",
		dashboards: dashboards,
		embed:      &cfg.Embed,
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
	targets := utils.MapArr(wt.dashboards, func(d dashboard) string { return d.id })
	return targets
}

func (wt *Watcher) Embed() *config.Embed {
	return wt.embed
}

func (wt *Watcher) CreateUrl(dashboard string) (string, error) {
	return fmt.Sprintf("%s/api/dashboards/%s/revisions", wt.baseUrl, dashboard), nil
}

func (wt *Watcher) CreateHref(target string, version string) *markdown.Href {
	text := fmt.Sprintf("%s revision %s", wt.getDashboardNameById(target), version)
	link := fmt.Sprintf("https://grafana.com/grafana/dashboards/%s/?tab=revisions", target)
	href := markdown.NewHref(text, link)
	return href
}

func (wt *Watcher) GetLatestVersion(data []byte, target string) (string, error) {
	var revisions Revisions

	if err := json.Unmarshal(data, &revisions); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	latest := revisions.Items[len(revisions.Items)-1]
	return strconv.Itoa(latest.Revision), nil
}

func (wt *Watcher) getDashboardNameById(id string) string {
	for _, d := range wt.dashboards {
		if d.id == id {
			return d.name
		}
	}
	// this case is not possible in general
	// therefore error is useless here
	return id
}
