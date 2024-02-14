package main

import (
	"log/slog"
	"os"
)

func NewLogger(lv slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lv}))
}
