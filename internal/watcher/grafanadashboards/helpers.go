package grafanadashboards

import (
	"fmt"
	"strconv"

	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

func createHrefs(vrs core.VersionRecords, dasdashboards []dashboard) []fmt.Stringer {
	hrefs := make([]fmt.Stringer, 0)
	for t, v := range vrs {
		text := fmt.Sprintf("%s revision %s", getDashboardNameById(t, dasdashboards), v)
		link := fmt.Sprintf("https://grafana.com/grafana/dashboards/%s/?tab=revisions", t)
		href := markdown.NewHref(text, link)
		hrefs = append(hrefs, href)
	}
	return hrefs
}

func getDashboardNameById(id string, dashboards []dashboard) string {
	for _, d := range dashboards {
		if d.id == id {
			return d.name
		}
	}
	// this case is not possible in general
	// therefore error is useless here
	return id
}

func getLatestRevision(revisions Revisions) string {
	latest := revisions.Items[len(revisions.Items)-1]
	return strconv.Itoa(latest.Revision)
}
