package dockerregistry

type Watcher struct {
	baseUrl   string
	authToken string
}

func NewWatcher() *Watcher {
	baseUrl := "https://hub.docker.com/v2"
	return &Watcher{baseUrl: baseUrl, authToken: ""}
}

func (w *Watcher) GetLastVersion() (string, error) {
	return "", nil
}
