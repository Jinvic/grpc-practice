package logger

import (
    "context"
    "log/slog"
)

type ctxKey struct{}

// WithContext 将 logger 注入到 context 中
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext 从 context 中获取 logger，如果不存在则返回全局默认
func FromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
        return logger
    }
    return defaultLogger
}

// 带 context 的日志方法（自动携带请求 ID 等信息）
func DebugCtx(ctx context.Context, msg string, args ...any) {
    FromContext(ctx).Debug(msg, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...any) {
    FromContext(ctx).Info(msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
    FromContext(ctx).Error(msg, args...)
}