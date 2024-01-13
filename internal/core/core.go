package core

import (
	"fmt"
	"os"

	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

func Task(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) {
	for _, wt := range wts {
		if err := task(wt, whs, r); err != nil {
			fmt.Fprintf(os.Stderr, "could not tick: %v", err)
		}
	}
}

func task(wt watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) error {
	targets := wt.GetTargets()
	versionRecords, err := watcher.GetLatestVersions(wt, targets)

	if err != nil {
		return fmt.Errorf("could not get version records %s: %v", wt.GetName(), err)
	}

	// update version record in database; notify on success
	for target, version := range versionRecords {
		updated, err := r.UpdateVersionRecord(target, version)
		if err != nil {
			return fmt.Errorf("could not update version record: %v", err)
		}
		// TODO: remove ! sign
		if !updated {
			for _, wh := range whs {
				title := wt.GetName()
				// TODO: add CreateVersionsUrl to Watcher
				// TODO: add GetTargetName to Watcher
				href := markdown.NewHref("text", "https://hub.docker.com/r/grafana/grafana/tags")
				if err := webhook.Notify(wh, title, href); err != nil {
					return fmt.Errorf("could not notify: %v", err)
				}
			}
		}
	}

	return nil
}
