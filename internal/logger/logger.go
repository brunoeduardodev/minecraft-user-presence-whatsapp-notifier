package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

var l *Logger

func Init() {
	l = &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}
}

func Info(ctx context.Context, msg string, keysAndValues ...any) {
	l.logger.InfoContext(ctx, msg, keysAndValues...)
}

func Error(ctx context.Context, err error, msg string, keysAndValues ...any) {
	l.logger.With("error", err).ErrorContext(ctx, msg, keysAndValues...)
}

func Warn(ctx context.Context, msg string, keysAndValues ...any) {
	l.logger.WarnContext(ctx, msg, keysAndValues...)
}

func Debug(ctx context.Context, msg string, keysAndValues ...any) {
	l.logger.DebugContext(ctx, msg, keysAndValues...)
}
