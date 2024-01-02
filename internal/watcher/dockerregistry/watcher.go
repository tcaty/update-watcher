package dockerregistry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
)

type Watcher struct {
	enabled bool
	name    string
	baseUrl string
	auth    auth
}

type auth struct {
	token    string
	login    string
	password string
}

func NewWatcher(cfg config.Dockerregistry) *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	return &Watcher{
		enabled: cfg.Enabled,
		name:    cfg.Name,
		baseUrl: baseUrl,
		auth: auth{
			token:    "",
			login:    cfg.Auth.Login,
			password: cfg.Auth.Password,
		},
	}
}

func (w *Watcher) IsEnabled() bool {
	return w.enabled
}

func (w *Watcher) GetName() string {
	return w.name
}

func (w *Watcher) Initialize() error {
	authUrl := fmt.Sprintf("%s/users/login", w.baseUrl)
	authData := map[string]string{"username": w.auth.login, "password": w.auth.password}
	jsonData, err := json.Marshal(authData)
	if err != nil {
		return err
	}

	resp, err := http.Post(authUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res map[string]string
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	w.auth.token = res["token"]

	return nil
}

func (w *Watcher) GetLatestVersions() (watcher.Versions, error) {
	w.getLatestTag()
	return nil, nil
}

func (w *Watcher) getLatestTag() (string, error) {
	return "", nil
}
