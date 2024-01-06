package main

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
	"github.com/tcaty/update-watcher/pkg/utils"
)

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	utils.HandleFatal("could not parse config", err)

	repo, err := initRepo(cfg.Postgresql)
	utils.HandleFatal("could not initialize repo", err)
	defer repo.Close()

	// startWatchers(cfg.Watchers, repo)
	s, err := initScheduler()
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

func initScheduler() (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("could not create scheduler: %v", err)
	}
	s.NewJob(
		gocron.CronJob("* * * * * *", true),
		gocron.NewTask(
			func() {
				fmt.Println("hello")
			},
		),
	)
	return s, nil
}

func startWatchers(cfg config.Watchers, repo *repository.Repository) {
	gdw := grafanadashboards.NewWatcher(cfg.Grafanadasboards)
	drw := dockerregistry.NewWatcher(cfg.Dockerregistry)

	watcher.Initialize[*grafanadashboards.Revisions](gdw)
	watcher.Initialize[*dockerregistry.Tags](drw)

	watcher.Tick[*grafanadashboards.Revisions](gdw, repo)
	watcher.Tick[*dockerregistry.Tags](drw, repo)
}
