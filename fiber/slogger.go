package main

import (
	"log/slog"
	"os"
)

var slogger *slog.Logger

func InitSlog() {
	slogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(slogger)
}
