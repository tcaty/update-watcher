package webhook

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/tcaty/update-watcher/pkg/utils"
)

type Webhook interface {
	IsEnabled() bool
	GetName() string
	GetUrl() string
	CreatePayload(target string, version string) (*bytes.Buffer, error)
}

func Notify(w Webhook, target string, version string) error {
	url := w.GetUrl()
	payload, err := w.CreatePayload(target, version)
	if err != nil {
		return fmt.Errorf("could not create http request empty payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return fmt.Errorf("http post request err: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %v", err)
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}

func Ping(w Webhook) error {
	url := w.GetUrl()
	empty, err := utils.CreateHttpRequestPayload(nil)
	if err != nil {
		return fmt.Errorf("could not create http request empty payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", empty)
	if err != nil {
		return fmt.Errorf("http post request err: %v", err)
	}

	// we expect that webhook with right url and empty payload
	// will send reponse with status code 400
	// in wrong url case it should send response with status code 401
	if resp.StatusCode != http.StatusBadRequest {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %v", err)
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}
