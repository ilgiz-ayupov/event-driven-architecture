package logger

import (
	"log/slog"
	"os"

	"event-driven-architecture/internal/usecase"
	"event-driven-architecture/pkg/slogconsole"
)

type slogLogger struct {
	log *slog.Logger
}

func NewSlogLogger() usecase.Logger {
	handler := slogconsole.New(
		os.Stderr,
		slogconsole.WithLevel(slog.LevelDebug),
		slogconsole.WithSource(true),
	)

	log := slog.New(handler)

	return &slogLogger{
		log: log,
	}
}

func (sl slogLogger) Debug(msg string, args ...any) {
	sl.log.Debug(msg, args...)
}

func (sl slogLogger) Info(msg string, args ...any) {
	sl.log.Info(msg, args...)
}

func (sl slogLogger) Warn(msg string, args ...any) {
	sl.log.Warn(msg, args...)
}

func (sl slogLogger) Error(msg string, args ...any) {
	sl.log.Error(msg, args...)
}
