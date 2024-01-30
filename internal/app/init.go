package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
	"github.com/tcaty/update-watcher/internal/webhook/discrod"
	"github.com/tcaty/update-watcher/pkg/inits"
)

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

func initRepo(cfg config.Postgresql) (core.Repository, error) {
	repo, err := repository.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := repo.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %v", err)
	}

	return repo, nil
}

func initWatchers(cfg config.Watchers) ([]core.Watcher, error) {
	wts := []core.Watcher{
		grafanadashboards.NewWatcher(cfg.Grafanadasboards),
		dockerregistry.NewWatcher(cfg.Dockerregistry),
	}
	wts = inits.ExcludeDisabledGeneric(wts)
	return wts, nil
}

func initWebhooks(cfg config.Webhooks) ([]core.Webhook, error) {
	whs := []core.Webhook{
		discrod.NewWebhook(cfg.Discord),
	}
	whs = inits.ExcludeDisabledGeneric(whs)
	if err := inits.PingAllGeneric(whs); err != nil {
		return nil, fmt.Errorf("unable ping webhook: %v", err)
	}
	return whs, nil
}

func initScheduler(cfg config.CronJob, core *core.Core) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("could not create scheduler: %v", err)
	}
	s.NewJob(
		gocron.CronJob(cfg.Crontab, cfg.WithSeconds),
		gocron.NewTask(core.WatchForUpdates),
	)
	return s, nil
}
