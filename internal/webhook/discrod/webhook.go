package discrod

import (
	"fmt"
	"net/http"

	"github.com/imroc/req/v3"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/core"
)

type Webhook struct {
	enabled bool
	name    string
	url     string
}

func NewWebhook(cfg config.Discord) *Webhook {
	return &Webhook{
		enabled: cfg.Enabled,
		name:    cfg.Name,
		url:     cfg.Url,
	}
}

func (wh *Webhook) Enabled() bool {
	return wh.enabled
}

func (wh *Webhook) String() string {
	return wh.name
}

func (wh *Webhook) Notify(msg core.Message) error {
	payload := createPayload(msg)

	_, err := req.C().R().
		SetBody(payload).
		Post(wh.url)

	return err
}

func (wh *Webhook) Ping() error {
	resp, err := req.C().R().
		Get(wh.url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response with status code %d", resp.StatusCode)
	}

	return nil
}
