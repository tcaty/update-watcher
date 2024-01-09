package grafanadashboards

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tcaty/update-watcher/internal/config"
)

type Watcher struct {
	enabled    bool
	name       string
	baseUrl    string
	dashboards []string
}

func NewWatcher(cfg config.Grafanadasboards) *Watcher {
	return &Watcher{
		enabled:    cfg.Enabled,
		name:       cfg.Name,
		baseUrl:    "https://grafana.com",
		dashboards: cfg.Dashboards,
	}
}

func (w *Watcher) IsEnabled() bool {
	return w.enabled
}

func (w *Watcher) GetName() string {
	return w.name
}

func (w *Watcher) GetTargets() []string {
	return w.dashboards
}

func (w *Watcher) CreateUrl(dashboard string) (string, error) {
	return fmt.Sprintf("%s/api/dashboards/%s/revisions", w.baseUrl, dashboard), nil
}

func (w *Watcher) GetLatestVersion(data []byte) (string, error) {
	var revisions Revisions

	if err := json.Unmarshal(data, &revisions); err != nil {
		return "", fmt.Errorf("cannot unmarshal json: %v", err)
	}

	latest := revisions.Items[len(revisions.Items)-1]
	return strconv.Itoa(latest.Revision), nil
}
