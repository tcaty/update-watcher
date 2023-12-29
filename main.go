package main

import (
	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
)

func main() {
	gdw := grafanadashboards.NewWatcher()
	drw := dockerregistry.NewWatcher()
	watcher.Tick(gdw)
	watcher.Tick(drw)
	cmd.Execute()
}
