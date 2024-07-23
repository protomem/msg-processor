package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger() *slog.Logger {
	o := tint.Options{Level: slog.LevelDebug}
	h := tint.NewHandler(os.Stdout, &o)
	return slog.New(h)
}
