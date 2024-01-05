package main

import (
	"fmt"
	"os"

	"github.com/tcaty/update-watcher/cmd"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/repository"
	"github.com/tcaty/update-watcher/internal/watcher"
	"github.com/tcaty/update-watcher/internal/watcher/dockerregistry"
	"github.com/tcaty/update-watcher/internal/watcher/grafanadashboards"
)

func main() {
	flags := cmd.Execute()
	cfg, err := config.Parse(flags.CfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse config: %v", err)
		os.Exit(1)
	}

	repo, err := initRepo(cfg.Postgresql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize repo: %v", err)
		os.Exit(1)
	}
	defer repo.Close()

	startWatchers(cfg.Watchers, repo)
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

func startWatchers(cfg config.Watchers, repo *repository.Repository) {
	gdw := grafanadashboards.NewWatcher(cfg.Grafanadasboards)
	drw := dockerregistry.NewWatcher(cfg.Dockerregistry)

	watcher.Initialize[*grafanadashboards.Revisions](gdw)
	watcher.Initialize[*dockerregistry.Tags](drw)

	watcher.Tick[*grafanadashboards.Revisions](gdw, repo)
	watcher.Tick[*dockerregistry.Tags](drw, repo)
}
