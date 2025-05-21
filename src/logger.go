package main

import (
	"log/slog"
	"os"
)

func NewLogger(isDebug bool) *slog.Logger {
	var level slog.Level
	if isDebug {
		level = slog.LevelDebug.Level()
	} else {
		level = slog.LevelInfo.Level()
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
