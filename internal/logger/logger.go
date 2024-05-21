package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

// Initialize устанавливает уровень логирования и создает новый логгер
func Initialize(logLevel string) {
	var handler slog.Handler

	switch logLevel {
	case "debug":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
	case "info":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case "warn":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})
	case "error":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError, AddSource: true})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	log = slog.New(handler)
}

// Logger возвращает текущий логгер
func Logger() *slog.Logger {
	return log
}
