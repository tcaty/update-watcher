package discrod

import (
	"bytes"
	"fmt"

	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Discord struct {
	enabled bool
	url     string
	avatar  string
	author  string
}

func NewWebhook(cfg config.Discord) *Discord {
	return &Discord{
		enabled: cfg.Enabled,
		url:     cfg.Url,
		avatar:  cfg.Avatar,
		author:  cfg.Author,
	}
}

func (w *Discord) IsEnabled() bool {
	return w.enabled
}

func (w *Discord) GetUrl() string {
	return w.url
}

func (w *Discord) CreatePayload(target string, version string) (*bytes.Buffer, error) {
	author := Author{
		Name:    "Update watcher",
		IconUrl: "https://cdn-icons-png.flaticon.com/512/1472/1472457.png",
	}
	embed := Embed{
		Author:      author,
		Title:       "New version detected",
		Description: fmt.Sprintf("New version here! %s:%s", target, version),
		Color:       "242424",
	}
	message := Message{
		Embeds: []Embed{embed},
	}
	payload, err := utils.CreateHttpRequestPayload(message)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
