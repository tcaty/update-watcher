package core

import (
	"fmt"

	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/webhook"
)

func Task(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) error {
	for _, wt := range wts {
		if err := task(wt, whs, r); err != nil {
			return fmt.Errorf("could not tick: %v", err)
		}
	}

	return nil
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
		if updated {
			if err := webhook.NotifyAll(whs, target, version); err != nil {
				return fmt.Errorf("could not notify all: %v", err)
			}
		}
	}

	return nil
}
