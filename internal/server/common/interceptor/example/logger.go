package exampleinterceptor

import (
	"bookstore/internal/pkg/logger"
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		l := logger.FromContext(ctx).With(slog.String("method", info.FullMethod))
		ctx = logger.WithContext(ctx, l)

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := status.Code(err)

		l.Info("request handled",
			"duration", duration.Milliseconds(),
			"status_code", statusCode.String(),
			"error", err,
		)

		if l.Enabled(ctx, slog.LevelDebug) {
			l.Debug("request details",
				"request", req,
				"response", resp,
			)
		}

		return resp, err
	}
}

func LoggingStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		l := logger.FromContext(stream.Context()).With(slog.String("method", info.FullMethod))

		wrappedStream := &loggedServerStream{
			ServerStream: stream,
			logger:       l,
		}

		err := handler(srv, wrappedStream)

		duration := time.Since(start)
		statusCode := status.Code(err)

		l.Info("request handled",
			"duration", duration.Milliseconds(),
			"status_code", statusCode.String(),
			"error", err,
			"recv_count", wrappedStream.recvCount,
			"send_count", wrappedStream.sendCount,
		)

		return err
	}
}

type loggedServerStream struct {
	grpc.ServerStream
	logger    *slog.Logger
	recvCount int
	sendCount int
}

func (s *loggedServerStream) RecvMsg(m any) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.recvCount++
		s.logger.Debug("received message",
			"count", s.recvCount,
			"message", m,
		)
	}
	return err
}

func (s *loggedServerStream) SendMsg(m any) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sendCount++
		s.logger.Debug("sent message",
			"count", s.sendCount,
			"message", m,
		)
	}
	return err
}

func (s *loggedServerStream) Context() context.Context {
	ctx := s.ServerStream.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	return logger.WithContext(ctx, s.logger)
}
