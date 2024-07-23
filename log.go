package main

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	o := slog.HandlerOptions{Level: slog.LevelDebug}
	h := slog.NewJSONHandler(os.Stdout, &o)
	return slog.New(h)
}
