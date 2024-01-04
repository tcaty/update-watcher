package main

import (
	"fmt"
	"os"

	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
)

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse config: %v", err)
		os.Exit(1)
	}

	repo, err := repository.New(cfg.Postgresql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer repo.Close()

	if false {
		startWatchers(cfg.Watchers)
	}
}

func startWatchers(cfg config.Watchers) {
	gdw := grafanadashboards.NewWatcher(cfg.Grafanadasboards)
	drw := dockerregistry.NewWatcher(cfg.Dockerregistry)

	watcher.Initialize[*grafanadashboards.Revisions](gdw)
	watcher.Initialize[*dockerregistry.Tags](drw)

	watcher.Tick[*grafanadashboards.Revisions](gdw)
	watcher.Tick[*dockerregistry.Tags](drw)
}
