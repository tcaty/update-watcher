package grafanadashboards

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Watcher struct {
	baseUrl string
}

func NewWatcher() *Watcher {
	return &Watcher{baseUrl: "https://grafana.com"}
}

func (w *Watcher) GetLastVersion() (string, error) {
	id := 1860
	url := w.getRevisionsUrl(id)
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

func (w *Watcher) getRevisionsUrl(id int) string {
	return fmt.Sprintf("%s/api/dashboards/%d/revisions", w.baseUrl, id)
}
