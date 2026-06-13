package logging

import (
	"io"
	"log/slog"
	"strings"
)

func New(level string, out io.Writer) *slog.Logger {
	var parsed slog.Level
	switch strings.ToLower(level) {
	case "debug":
		parsed = slog.LevelDebug
	case "warn", "warning":
		parsed = slog.LevelWarn
	case "error":
		parsed = slog.LevelError
	default:
		parsed = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: parsed}))
}
