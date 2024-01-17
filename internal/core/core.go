package core

import (
	"fmt"

	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

func WatchForUpdates(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) {
	for _, wt := range wts {
		if err := watchForUpdates(wt, whs, r); err != nil {
			wt.Slog().Error("could not watch for updates", "error", err)
		}
	}
}

func watchForUpdates(wt watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) error {
	updatedTargetsHrefs := make([]*markdown.Href, 0)
	targets := wt.Targets()
	versionRecords, err := watcher.FetchLatestVersionRecords(wt, targets)

	if err != nil {
		return fmt.Errorf("could not fetch latest version records %s: %v", wt.Name(), err)
	}

	for t, v := range versionRecords {
		updated, err := r.UpdateVersionRecord(t, v)
		if err != nil {
			return fmt.Errorf("could not update version record: %v", err)
		}
		// TODO: remove ! sign
		if !updated {
			href := wt.CreateHref(t, v)
			updatedTargetsHrefs = append(updatedTargetsHrefs, href)
		}
	}

	for _, wh := range whs {
		title := wt.Name()
		if err := webhook.Notify(wh, title, updatedTargetsHrefs); err != nil {
			return fmt.Errorf("could not notify: %v", err)
		}
	}

	return nil
}
