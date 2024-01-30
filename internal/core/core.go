package core

import (
	"fmt"

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
		vrs, err := wt.FetchLatestVersionRecords()
		if err != nil {
			fmt.Println(err)
		}

		updatedVrs := c.updateVersionRecords(vrs)

		if len(updatedVrs) == 0 {
			return
		}

		msg := wt.CreateMessageAboutUpdates(updatedVrs)
		c.notifyAboutUpdates(msg)
	}
}

func (c *Core) updateVersionRecords(vrs []entities.VersionRecord) []entities.VersionRecord {
	updatedVrs := make([]entities.VersionRecord, 0)

	for _, vr := range vrs {
		updated, err := c.repo.UpdateVersionRecord(vr)
		if err != nil {
			fmt.Println(err)
		}

		if updated {
			updatedVrs = append(updatedVrs, vr)
		}
	}

	return updatedVrs
}

func (c *Core) notifyAboutUpdates(msg Message) {
	for _, wh := range c.whs {
		if err := wh.Notify(msg); err != nil {
			fmt.Println(err)
		}
	}
}
