package discrod

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Webhook struct {
	slog    *slog.Logger
	enabled bool
	name    string
	url     string
}

func NewWebhook(cfg config.Discord) *Webhook {
	return &Webhook{
		slog:    slog.Default().With("webhook", cfg.Name),
		enabled: cfg.Enabled,
		name:    cfg.Name,
		url:     cfg.Url,
	}
}

func (wh *Webhook) Slog() *slog.Logger {
	return wh.slog
}

func (wh *Webhook) Enabled() bool {
	return wh.enabled
}

func (wh *Webhook) Name() string {
	return wh.name
}

func (wh *Webhook) Url() string {
	return wh.url
}

func (wh *Webhook) CreatePayload(msg *webhook.Message) (*bytes.Buffer, error) {
	author := Author{
		Name:    msg.Author,
		IconUrl: msg.Avatar,
	}
	embed := Embed{
		Author:      author,
		Description: msg.Description,
		Color:       msg.Color,
	}
	message := Message{
		Embeds: []Embed{embed},
	}
	payload, err := utils.CreateHttpRequestPayload(message)
	if err != nil {
		return nil, fmt.Errorf("could not create http request payload: %v", err)
	}
	return payload, nil
}
