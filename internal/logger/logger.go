package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Setup(level string) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLevel(level),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func getLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
