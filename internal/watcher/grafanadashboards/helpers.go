package grafanadashboards

import (
	"fmt"
	"strconv"

	"github.com/tcaty/update-watcher/internal/entities"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

func createHrefs(vrs []entities.VersionRecord, dasdashboards []dashboard) []fmt.Stringer {
	hrefs := make([]fmt.Stringer, 0)
	for _, vr := range vrs {
		name := getDashboardNameById(vr.Target, dasdashboards)
		text := fmt.Sprintf("%s revision %s", name, vr.Version)
		link := fmt.Sprintf("https://grafana.com/grafana/dashboards/%s/?tab=revisions", vr.Target)
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
