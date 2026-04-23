package logger

import (
	"bookstore/internal/pkg/config"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

var defaultLogger *slog.Logger

func InitLogger(cfg *config.Logging, logFile string) {
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
		writer = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			Compress:   cfg.Compress,
			LocalTime:  cfg.LocalTime,
		}
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
}

func GetLogger() *slog.Logger {
	return defaultLogger
}

func Debug(msg string, args ...any) { defaultLogger.Debug(msg, args...) }
func Info(msg string, args ...any)  { defaultLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { defaultLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { defaultLogger.Error(msg, args...) }
