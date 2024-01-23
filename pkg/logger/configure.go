package logger

import (
	"log/slog"
	"os"
)

// Configure the logger
func Configure(level, format string) {
	logLevel := slog.LevelDebug
	switch level {
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	var logger *slog.Logger
	if format == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &opts))
	}

	slog.SetDefault(logger)
}
