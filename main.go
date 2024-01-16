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
// TODO: update message format

func main() {
	initLogger()

	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	utils.HandleFatal("could not parse config", err)

	slog.Info("Initializing repo...")
	repo, err := initRepo(cfg.Postgresql)
	utils.HandleFatal("could not initialize repo", err)
	defer repo.Close()
	slog.Info("Repo initialized successfully")

	slog.Info("Initializing watchers...")
	wts, err := initWatchers(cfg.Watchers)
	utils.HandleFatal("could not initialize watchers", err)
	slog.Info("Watchers initialized successfully")

	slog.Info("Initializing webhooks...")
	whs, err := initWebhooks(cfg.Webhooks)
	utils.HandleFatal("could not initialize webhooks", err)
	slog.Info("Webhooks initialized successfully")

	slog.Info("Initializing scheduler...")
	s, err := initScheduler(cfg.CronJob, wts, whs, repo)
	utils.HandleFatal("could not initialize scheduler", err)
	// s.Start()
	defer s.Shutdown()
	slog.Info("Repo initialized successfully")

	// block current channel to run cronjob
	// see: https://github.com/go-co-op/gocron/issues/647
	select {}
}

func initLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func initRepo(cfg config.Postgresql) (*repository.Repository, error) {
	repo, err := repository.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err := repo.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}
	if err := repo.InitializeTables(); err != nil {
		return nil, fmt.Errorf("unable to initialize database tables: %v", err)
	}
	return repo, nil
}

func initWatchers(cfg config.Watchers) ([]watcher.Watcher, error) {
	watchers := ([]watcher.Watcher{
		grafanadashboards.NewWatcher(cfg.Grafanadasboards),
		dockerregistry.NewWatcher(cfg.Dockerregistry),
	})
	return watchers, nil
}

func initWebhooks(cfg config.Webhooks) ([]webhook.Webhook, error) {
	webhooks := []webhook.Webhook{
		discrod.NewWebhook(cfg.Discord),
	}
	for _, w := range webhooks {
		if err := webhook.Ping(w); err != nil {
			return nil, fmt.Errorf("could not ping webhook %s: %v", w.Name(), err)
		}
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
		gocron.NewTask(core.Task, wts, whs, r),
	)
	return s, nil
}
