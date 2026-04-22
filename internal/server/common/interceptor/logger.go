package interceptor

import (
	logger "bookstore/internal/pkg/logger"
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

var (
	interceptorLogger *slog.Logger
)

func InitLoggingInterceptor() error {
	interceptorLogger = logger.GetLogger()
	return nil
}

func InterceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		interceptorLogger.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func LoggingUnaryInterceptor(opts ...logging.Option) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(InterceptorLogger(), opts...)
}

func LoggingStreamInterceptor(opts ...logging.Option) grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(InterceptorLogger(), opts...)
}
