package discrod

import "github.com/tcaty/update-watcher/internal/core"

func createPayload(msg core.Message) Payload {
	author := Author{
		Name:    msg.Author,
		IconUrl: msg.Avatar,
	}
	embed := Embed{
		Author:      author,
		Description: msg.Description,
		Color:       msg.Color,
	}
	message := Payload{
		Embeds: []Embed{embed},
	}

	return message
}
