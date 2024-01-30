package main

import (
	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/app"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/utils"
)

// TODO move logger to pkg/, refactor logging implementation
// TODO: move config files to config/, refactor this with cleanenv
// TODO: refactor watchers
// TODO: implement gorountines
// TODO: add Makefile
// TODO: update README

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)

	if err != nil {
		utils.HandleFatal("could not parse config", err)
	}

	app.Run(*cfg)
}
