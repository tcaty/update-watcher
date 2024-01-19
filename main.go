package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/internal/webhook/discrod"
	"github.com/tcaty/update-watcher/pkg/utils"
)

// TODO: complete logging
// TODO: set default values

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	if err != nil {
		utils.HandleFatal("could not parse config", err)
	}

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

	slog.Info("initializing scheduler...")
	s, err := initScheduler(cfg.CronJob, wts, whs, repo)
	if err != nil {
		utils.HandleFatal("could not initialize scheduler", err)
	}
	s.Start()
	defer s.Shutdown()

	slog.Info("everything is ready. starting watching for updates.")
	// block current channel to run cronjob
	// see: https://github.com/go-co-op/gocron/issues/647
	select {}
}

func initLogger(cfg config.Logger) error {
	level := new(slog.LevelVar)
	err := level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return fmt.Errorf("unable to parse logLevel")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
	return nil
}

func initRepo(cfg config.Postgresql) (*repository.Repository, error) {
	repo, err := repository.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	repo.Slog().Debug("connection established")

	if err := repo.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}
	repo.Slog().Debug("ping is successful")

	if err := repo.InitializeTables(); err != nil {
		return nil, fmt.Errorf("unable to initialize database tables: %v", err)
	}
	repo.Slog().Debug("tables initialized successfully")

	return repo, nil
}

func initWatchers(cfg config.Watchers) ([]watcher.Watcher, error) {
	watchers := []watcher.Watcher{
		grafanadashboards.NewWatcher(cfg.Grafanadasboards),
		dockerregistry.NewWatcher(cfg.Dockerregistry),
	}
	return watchers, nil
}

func initWebhooks(cfg config.Webhooks) ([]webhook.Webhook, error) {
	webhooks := []webhook.Webhook{
		discrod.NewWebhook(cfg.Discord),
	}
	for _, wh := range webhooks {
		if err := webhook.Ping(wh); err != nil {
			return nil, fmt.Errorf("could not ping webhook %s: %v", wh.Name(), err)
		}
		wh.Slog().Debug("ping is successful")
	}
	return webhooks, nil
}

func initScheduler(cfg config.CronJob, wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("could not create scheduler: %v", err)
	}
	s.NewJob(
		gocron.CronJob(cfg.Crontab, cfg.WithSeconds),
		gocron.NewTask(core.WatchForUpdates, wts, whs, r),
	)
	return s, nil
}
