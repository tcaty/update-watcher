package core

import (
	"fmt"
	"os"

	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/webhook"
)

func Task(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) {
	for _, wt := range wts {
		if err := task(wt, whs, r); err != nil {
			fmt.Fprintf(os.Stderr, "could not tick: %v", err)
		}
	}
}

func task(wt watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) error {
	targets := wt.Targets()
	versionRecords, err := watcher.GetLatestVersions(wt, targets)

	if err != nil {
		return fmt.Errorf("could not get version records %s: %v", wt.Name(), err)
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
				title := wt.Name()
				href := wt.CreateHref(target, version)
				if err := webhook.Notify(wh, title, href); err != nil {
					return fmt.Errorf("could not notify: %v", err)
				}
			}
		}
	}

	return nil
}
