package core

import "fmt"

type Repository interface {
	UpdateVersionRecord(target string, version string) (bool, error)
}

type Watcher interface {
	FetchLatestVersionRecords() (VersionRecords, error)
	CreateMessageAboutUpdates(vrs VersionRecords) Message
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

func (core *Core) WatchForUpdates() {
	for _, wt := range core.wts {
		vrs, err := wt.FetchLatestVersionRecords()
		if err != nil {
			fmt.Println("error")
		}

		updatedVrs := core.updateVersionRecords(vrs)

		if len(updatedVrs) == 0 {
			return
		}

		msg := wt.CreateMessageAboutUpdates(updatedVrs)
		core.notifyAboutUpdates(msg)
	}
}

func (core *Core) updateVersionRecords(vrs VersionRecords) VersionRecords {
	updatedVrs := make(VersionRecords)

	for t, v := range vrs {
		updated, err := core.repo.UpdateVersionRecord(t, v)
		if err != nil {
			fmt.Println("error")
		}

		if updated {
			updatedVrs[t] = v
		}
	}

	return updatedVrs
}

func (core *Core) notifyAboutUpdates(msg Message) {
	for _, wh := range core.whs {
		if err := wh.Notify(msg); err != nil {
			fmt.Println("error")
		}
	}
}
