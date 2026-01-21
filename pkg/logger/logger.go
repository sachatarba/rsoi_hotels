package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var logHandler slog.Handler

	opts := &slog.HandlerOptions{
		Level: func() slog.Level {
			if env == "local" || env == "dev" {
				return slog.LevelDebug
			}
			return slog.LevelInfo
		}(),
		AddSource: env == "local",
	}

	if env == "local" {
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(logHandler)

	return logger
}
