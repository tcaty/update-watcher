package utils

import (
	"log/slog"
	"os"
)

func HandleFatal(message string, err error) {
	slog.Error(message, "error", err)
	os.Exit(1)
}
