package main

import (
	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/app"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/utils"
)

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)

	if err != nil {
		utils.HandleFatal("could not parse config", err)
	}

	app.Run(*cfg)
}
