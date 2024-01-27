package app

import (
	"log/slog"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/pkg/utils"
)

func Run(cfg config.Config) {
	if err := initLogger(cfg.Logger); err != nil {
		utils.HandleFatal("could not init logger", err)
	}

	slog.Info("initializing repo...")
	repo, err := initRepo(cfg.Postgresql)
	if err != nil {
		utils.HandleFatal("could not initialize repo", err)
	}
	defer repo.Close()

	slog.Info("initializing watchers...")
	wts, err := initWatchers(cfg.Watchers)
	if err != nil {
		utils.HandleFatal("could not initialize watchers", err)
	}

	slog.Info("initializing webhooks...")
	whs, err := initWebhooks(cfg.Webhooks)
	if err != nil {
		utils.HandleFatal("could not initialize webhooks", err)
	}

	core := core.New(repo, wts, whs)

	slog.Info("initializing scheduler...")
	s, err := initScheduler(cfg.CronJob, core)
	if err != nil {
		utils.HandleFatal("could not initialize scheduler", err)
	}
	s.Start()
	defer s.Shutdown()

	slog.Info("everything is ready. starting watching for updates.")

	if cfg.CronJob.ExecImmediate {
		core.WatchForUpdates()
	}

	// block current channel to run cronjob
	// see: https://github.com/go-co-op/gocron/issues/647
	select {}
}
