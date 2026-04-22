package exampleinterceptor

import (
	"context"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	validator protovalidate.Validator
)

func InitValidateInterceptor() error {
	var err error
	validator, err = protovalidate.New()
	if err != nil {
		return err
	}
	return nil
}

func ValidateUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.Internal, "unsupported message type: %T", req)
		}
		if err := validator.Validate(msg); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
		}
		return handler(ctx, req)
	}
}

func ValidateStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &wrappedServerStream{
			ServerStream: ss,
			validator:    validator,
		})
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	validator protovalidate.Validator
}

func (w *wrappedServerStream) RecvMsg(m any) error {
	if err := w.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	msg, ok := m.(proto.Message)
	if !ok {
		return status.Errorf(codes.Internal, "unsupported message type: %T", m)
	}
	return w.validator.Validate(msg)
}
