package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func InitLogger(level string) {
	logLevel := &slog.LevelVar{}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	Log = slog.New(handler)

	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel.Set(slog.LevelDebug)
	case "INFO":
		logLevel.Set(slog.LevelInfo)
	case "WARN":
		logLevel.Set(slog.LevelWarn)
	case "ERROR":
		logLevel.Set(slog.LevelError)
	default:
		logLevel.Set(slog.LevelError)
	}
}
