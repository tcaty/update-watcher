package utils

import (
	"bytes"
	"encoding/json"
)

func CreateHttpRequestPayload(v any) (*bytes.Buffer, error) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(v)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
