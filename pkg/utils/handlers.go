package utils

import (
	"log/slog"
	"os"
)

func HandleFatal(message string, err error) {
	slog.Error(message, err)
	os.Exit(1)
}
