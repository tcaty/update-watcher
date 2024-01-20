package core

import (
	"fmt"
	"log/slog"

	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/pkg/markdown"
	"github.com/tcaty/update-watcher/pkg/utils"
)

func WatchForUpdates(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) {
	total, updated := 0, 0

	slog.Info("starting watching for updates cycle")

	for _, wt := range wts {

		wt.Slog().Debug("watching for updates")

		wtTotal := len(wt.Targets())
		wtUpdated, err := watchForUpdates(wt, whs, r)

		if err != nil {
			wt.Slog().Error("could not watch for updates", "error", err)
		}

		total += wtTotal
		updated += wtUpdated

		wt.Slog().Debug("processing results", "targets_total", wtTotal, "targets_updated", wtUpdated)
	}

	slog.Info("processing results", "targets_total", total, "targets_updated", updated)
}

func watchForUpdates(wt watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) (int, error) {
	updatedTargetsHrefs := make([]*markdown.Href, 0)
	targets := wt.Targets()
	versionRecords, err := watcher.FetchLatestVersionRecords(wt, targets)

	if err != nil {
		return 0, fmt.Errorf("could not fetch latest version records %s: %v", wt.Name(), err)
	}

	for t, v := range versionRecords {
		wt.Slog().Debug("starting version record update", "target", t, "version", v)
		updated, err := r.UpdateVersionRecord(t, v)

		if err != nil {
			// should we return err here or just write to log like this?
			wt.Slog().Error("could not update version record", "error", err)
		}

		if updated {
			href := wt.CreateHref(t, v)
			updatedTargetsHrefs = append(updatedTargetsHrefs, href)
		}

		wt.Slog().Debug("processing results", "target", t, "version", v, "updated", updated)
	}

	// if there are no relevant updates
	// then we don't need any notifications
	// stop func execution without error
	if len(updatedTargetsHrefs) == 0 {
		return 0, nil
	}

	whNotified := 0
	for _, wh := range whs {
		msg := createMessage(wt, updatedTargetsHrefs)
		if err := webhook.Notify(wh, msg); err != nil {
			// should we return err here or just write to log like this?
			wh.Slog().Error("could not notify", "error", err)
		}
		whNotified += 1
		wt.Slog().Debug(
			"webhook notified",
			"webhook", wh.Name(),
			"targets_total", len(targets),
			"targets_updated", len(updatedTargetsHrefs),
		)
	}
	wt.Slog().Info("webhooks notified", "notified_count", whNotified)

	return len(updatedTargetsHrefs), nil
}

func createMessage(wt watcher.Watcher, hrefs []*markdown.Href) *webhook.Message {
	list := markdown.CreateUnorderedList(
		utils.MapArr(hrefs, func(h *markdown.Href) string { return h.Sprint() }),
	)
	descr := fmt.Sprintf("%s\n%s", wt.Embed().Text, list)
	msg := &webhook.Message{
		Author:      wt.Name(),
		Avatar:      wt.Embed().Avatar,
		Description: descr,
		Color:       wt.Embed().Color,
	}
	return msg
}
