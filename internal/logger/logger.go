package logger

import (
	"log/slog"
	"os"
)

var lg Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource:   false,
	Level:       slog.LevelDebug,
	ReplaceAttr: nil,
}))

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
}

func Default() Logger {
	return lg
}

func Fatal(msg string, args ...any) {
	lg.Warn(msg, args...)
	os.Exit(1)
}
