package utils

import (
	"bytes"
	"encoding/json"
	"log/slog"
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
		slog.Error(message, err)
		os.Exit(1)
	}
}

func MapArr[T comparable, V comparable](arr []T, callback func(v T) V) []V {
	newArr := make([]V, len(arr))
	for i, v := range arr {
		newArr[i] = callback(v)
	}
	return newArr
}
