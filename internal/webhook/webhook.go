package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Webhook interface {
	IsEnabled() bool
	GetUrl() string
	CreatePayload(target string, version string) (*bytes.Buffer, error)
}

func NotifyAll(whs []Webhook, target string, version string) error {
	for _, w := range whs {
		if err := notify(w, target, version); err != nil {
			return fmt.Errorf("could not notify: %v", err)
		}
	}
	return nil
}

func notify(w Webhook, target string, version string) error {
	url := w.GetUrl()
	payload, err := w.CreatePayload(target, version)

	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}
