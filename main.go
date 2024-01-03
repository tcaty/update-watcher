package main

import (
	"fmt"

	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
)

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	if err != nil {
		fmt.Printf("could not parse config: %v", err)
		return
	}

	gdw := grafanadashboards.NewWatcher(cfg.Watchers.Grafanadasboards)
	drw := dockerregistry.NewWatcher(cfg.Watchers.Dockerregistry)

	watcher.Initialize[*grafanadashboards.Revisions](gdw)
	watcher.Initialize[*dockerregistry.Tags](drw)

	watcher.Tick[*grafanadashboards.Revisions](gdw)
	watcher.Tick[*dockerregistry.Tags](drw)
}
