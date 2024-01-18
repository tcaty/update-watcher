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

func (w *Watcher) Slog() *slog.Logger {
	return w.slog
}

func (w *Watcher) Enabled() bool {
	return w.enabled
}

func (w *Watcher) Name() string {
	return w.name
}

func (w *Watcher) Targets() []string {
	targets := utils.MapArr(w.dashboards, func(d dashboard) string { return d.id })
	return targets
}

func (w *Watcher) Embed() *config.Embed {
	return w.embed
}

func (w *Watcher) CreateUrl(dashboard string) (string, error) {
	return fmt.Sprintf("%s/api/dashboards/%s/revisions", w.baseUrl, dashboard), nil
}

func (w *Watcher) CreateHref(target string, version string) *markdown.Href {
	text := fmt.Sprintf("%s revision %s", w.getDashboardNameById(target), version)
	link := fmt.Sprintf("https://grafana.com/grafana/dashboards/%s/?tab=revisions", target)
	href := markdown.NewHref(text, link)
	return href
}

func (w *Watcher) GetLatestVersion(data []byte, target string) (string, error) {
	var revisions Revisions

	if err := json.Unmarshal(data, &revisions); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	latest := revisions.Items[len(revisions.Items)-1]
	return strconv.Itoa(latest.Revision), nil
}

func (w *Watcher) getDashboardNameById(id string) string {
	for _, d := range w.dashboards {
		if d.id == id {
			return d.name
		}
	}
	// this case is not possible in general
	// therefore error is useless here
	return id
}
