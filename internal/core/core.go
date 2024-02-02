package core

import (
	"log/slog"

	"github.com/tcaty/update-watcher/internal/entities"
)

type Repository interface {
	UpdateVersionRecord(vr entities.VersionRecord) (bool, error)
}

type Watcher interface {
	FetchLatestVersionRecords() ([]entities.VersionRecord, error)
	CreateMessageAboutUpdates(vrs []entities.VersionRecord) Message
}

type Webhook interface {
	Notify(msg Message) error
}

type Core struct {
	repo Repository
	wts  []Watcher
	whs  []Webhook
}

func New(repo Repository, wts []Watcher, whs []Webhook) *Core {
	return &Core{
		repo: repo,
		wts:  wts,
		whs:  whs,
	}
}

func (c *Core) WatchForUpdates() {
	for _, wt := range c.wts {
		log := slog.With("watcher", wt)

		log.Info("starting watching for updates cycle")

		vrs, err := wt.FetchLatestVersionRecords()

		if err != nil {
			log.Error("error occured while fetching version records", "error", err)
		}

		updatedVrs := c.updateVersionRecords(vrs)

		log.Info(
			"processing results",
			"targets_watched", len(vrs),
			"targets_updated", len(updatedVrs),
		)

		if len(updatedVrs) == 0 {
			log.Info("there are no updates. notifications won't be sent")
			continue
		}

		msg := wt.CreateMessageAboutUpdates(updatedVrs)
		c.notifyAboutUpdates(msg)
	}
}

func (c *Core) updateVersionRecords(vrs []entities.VersionRecord) []entities.VersionRecord {
	updatedVrs := make([]entities.VersionRecord, 0)

	for _, vr := range vrs {
		log := slog.With(
			"target", vr.Target,
			"version", vr.Version,
		)

		updated, err := c.repo.UpdateVersionRecord(vr)

		if err != nil {
			log.Error("error occured while updating version record", "error", err)
		}

		if updated {
			updatedVrs = append(updatedVrs, vr)
		}

		log.Debug("updating version record", "updated", updated)
	}

	return updatedVrs
}

func (c *Core) notifyAboutUpdates(msg Message) {
	for _, wh := range c.whs {
		log := slog.With(
			"webhook", wh,
			"watcher", msg.Author,
		)

		if err := wh.Notify(msg); err != nil {
			log.Error("error occured while notifying", "error", err)
		} else {
			log.Info("successfully notified")
		}
	}
}
