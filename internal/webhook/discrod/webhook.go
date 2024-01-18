package discrod

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/webhook"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Discord struct {
	slog    *slog.Logger
	enabled bool
	name    string
	url     string
}

func NewWebhook(cfg config.Discord) *Discord {
	return &Discord{
		slog: slog.Default().With("webhook", cfg.Name),
		name: cfg.Name,
		url:  cfg.Url,
	}
}

func (w *Discord) Slog() *slog.Logger {
	return w.slog
}

func (w *Discord) Enabled() bool {
	return w.enabled
}

func (w *Discord) Name() string {
	return w.name
}

func (w *Discord) Url() string {
	return w.url
}

func (w *Discord) CreatePayload(msg *webhook.Message) (*bytes.Buffer, error) {
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
