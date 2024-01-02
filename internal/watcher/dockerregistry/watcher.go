package dockerregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/watcher"
)

type Watcher struct {
	enabled bool
	name    string
	baseUrl string
	images  []string
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
		images:  cfg.Images,
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
	updates := make(watcher.Versions, len(w.images))
	for _, d := range w.images {
		r, err := w.getLatestTag(d)
		if err != nil {
			return nil, err
		}
		updates[d] = r
	}
	return updates, nil
}

func (w *Watcher) getLatestTag(image string) (string, error) {
	url, err := w.getRepositoryTagsUrl(image)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get tags url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read tags body: %v", err)
	}

	var tags Tags
	if err := json.Unmarshal(body, &tags); err != nil {
		return "", fmt.Errorf("cannot unmarshal tags: %v", err)
	}

	latestTag := tags.Results[1]
	return latestTag.Name, nil
}

func (w *Watcher) getRepositoryTagsUrl(image string) (string, error) {
	b := []byte(image)
	i := bytes.IndexByte(b, byte('/'))
	if i < 0 {
		return "", errors.New("docker image should fit the format {namespace}/{repository}")
	}
	ns, repo := string(b[:i]), string(b[i+1:])
	return fmt.Sprintf("%s/namespaces/%s/repositories/%s/tags", w.baseUrl, ns, repo), nil
}
