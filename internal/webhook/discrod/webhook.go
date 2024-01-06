package discrod

import (
	"bytes"

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
	image := Image{
		Url: "",
	}
	embed := Embed{
		Image: &image,
	}
	message := Message{
		Username: "",
		Content:  "",
		Embeds:   &[]Embed{embed},
	}
	payload, err := utils.CreateHttpRequestPayload(message)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
