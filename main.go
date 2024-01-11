package main

import (
	"fmt"

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

// TODO: refactor watcher
// TODO: add crontab to config
// TODO: create task for cronjob

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	utils.HandleFatal("could not parse config", err)

	repo, err := initRepo(cfg.Postgresql)
	utils.HandleFatal("could not initialize repo", err)
	defer repo.Close()

	wts, err := initWatchers(cfg.Watchers)
	utils.HandleFatal("could not initialize watchers", err)

	whs, err := initWebhooks(cfg.Webhooks)
	utils.HandleFatal("could not initialize webhooks", err)

	s, err := initScheduler(wts, whs, repo)
	utils.HandleFatal("could not initialize scheduler", err)
	s.Start()
	defer s.Shutdown()

	// block current channel to run cronjob
	// see: https://github.com/go-co-op/gocron/issues/647
	select {}
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
	watchers := []watcher.Watcher{
		grafanadashboards.NewWatcher(cfg.Grafanadasboards),
		dockerregistry.NewWatcher(cfg.Dockerregistry),
	}
	for _, w := range watchers {
		if err := watcher.Initialize(w); err != nil {
			return nil, fmt.Errorf("could not initialize watcher %s: %v", w.GetName(), err)
		}
	}
	return watchers, nil
}

func initWebhooks(cfg config.Webhooks) ([]webhook.Webhook, error) {
	webhooks := []webhook.Webhook{
		discrod.NewWebhook(cfg.Discord),
	}
	for _, w := range webhooks {
		if err := webhook.Ping(w); err != nil {
			return nil, fmt.Errorf("could not ping webhook %s: %v", w.GetName(), err)
		}
	}
	return webhooks, nil
}

func initScheduler(wts []watcher.Watcher, whs []webhook.Webhook, r *repository.Repository) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("could not create scheduler: %v", err)
	}
	s.NewJob(
		gocron.CronJob("* * * * * *", true),
		gocron.NewTask(core.Task, wts, whs, r),
	)
	return s, nil
}
