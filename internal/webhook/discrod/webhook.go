package discrod

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Discord struct {
	slog    *slog.Logger
	enabled bool
	name    string
	url     string
	avatar  string
	author  string
	color   int
}

func NewWebhook(cfg config.Discord) *Discord {
	return &Discord{
		slog:   slog.Default().With("webhook", cfg.Name),
		name:   cfg.Name,
		url:    cfg.Url,
		avatar: cfg.Avatar,
		author: cfg.Author,
		color:  cfg.Color,
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

func (w *Discord) CreatePayload(title string, description string) (*bytes.Buffer, error) {
	author := Author{
		Name:    w.author,
		IconUrl: w.avatar,
	}
	embed := Embed{
		Author:      author,
		Title:       title,
		Description: description,
		Color:       w.color,
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
