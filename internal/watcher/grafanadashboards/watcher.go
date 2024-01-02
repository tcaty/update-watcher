package grafanadashboards

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
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

func (w *Watcher) Initialize() error {
	return nil
}

func (w *Watcher) GetLatestVersions() (watcher.Versions, error) {
	updates := make(watcher.Versions, len(w.dashboards))
	for _, d := range w.dashboards {
		r, err := w.getLatestRevision(d)
		if err != nil {
			return nil, err
		}
		updates[d] = r
	}
	return updates, nil
}

func (w *Watcher) getLatestRevision(dashboard string) (string, error) {
	url := w.getRevisionsUrl(dashboard)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get grafanadashboards url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read revisions body: %v", err)
	}

	var revisions Revisions
	if err := json.Unmarshal(body, &revisions); err != nil {
		return "", fmt.Errorf("cannot unmarshal revisions: %v", err)
	}

	latestRevision := revisions.Items[len(revisions.Items)-1]
	return strconv.Itoa(latestRevision.Revision), nil
}

func (w *Watcher) getRevisionsUrl(dashboard string) string {
	return fmt.Sprintf("%s/api/dashboards/%s/revisions", w.baseUrl, dashboard)
}
