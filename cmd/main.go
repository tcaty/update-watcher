package main

import (
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
)

func main() {
	gdw := grafanadashboards.NewWatcher()
	watcher.Tick(gdw)
}
