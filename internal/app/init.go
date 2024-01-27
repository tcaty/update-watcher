package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
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
	filtered := utils.FilterArr(
		watchers,
		func(wt watcher.Watcher) bool {
			wt.Slog().Debug("filtering watchers", "enabled", wt.Enabled())
			return wt.Enabled()
		},
	)
	return filtered, nil
}

func initWebhooks(cfg config.Webhooks) ([]webhook.Webhook, error) {
	webhooks := []webhook.Webhook{
		discrod.NewWebhook(cfg.Discord),
	}
	filtered := utils.FilterArr(
		webhooks,
		func(wh webhook.Webhook) bool {
			wh.Slog().Debug("filtering webhooks", "enabled", wh.Enabled())
			return wh.Enabled()
		},
	)
	for _, wh := range filtered {
		if err := webhook.Ping(wh); err != nil {
			return nil, fmt.Errorf("could not ping webhook %s: %v", wh.Name(), err)
		}
		wh.Slog().Debug("ping is successful")
	}
	return filtered, nil
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
