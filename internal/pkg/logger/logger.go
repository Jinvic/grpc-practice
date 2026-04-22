package logger

import (
	"bookstore/internal/pkg/config"
	fileUtil "bookstore/util/file"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var defaultLogger *slog.Logger

func InitLogger(cfg *config.Logging) (cleanup func(), err error) {
	cleanup = func() {}

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	var writer io.Writer
	switch cfg.Output {
	case "stdout":
		writer = os.Stdout
	case "file":
		if err := fileUtil.MkDir(filepath.Dir(cfg.File)); err != nil {
			return cleanup, fmt.Errorf("failed to create directory: %w", err)
		}
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return cleanup, fmt.Errorf("failed to open file: %w", err)
		}
		cleanup = func() {
			file.Close()
		}

		writer = file
	default:
		writer = os.Stdout
	}

	var handler slog.Handler
	options := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(writer, options)
	case "text":
		handler = slog.NewTextHandler(writer, options)
	default:
		handler = slog.NewTextHandler(writer, options)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
	return cleanup, nil
}

func GetLogger() *slog.Logger {
	return defaultLogger
}

func Debug(msg string, args ...any) { defaultLogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { defaultLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { defaultLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { defaultLogger.Error(msg, args...) }
