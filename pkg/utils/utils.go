package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func CreateHttpRequestPayload(v any) (*bytes.Buffer, error) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(v)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func HandleFatal(message string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v", message, err)
		os.Exit(1)
	}
}
