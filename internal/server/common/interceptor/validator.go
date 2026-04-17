package interceptor

import (
	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
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

func ValidateUnaryInterceptor(opts ...protovalidate_middleware.Option) grpc.UnaryServerInterceptor {
	return protovalidate_middleware.UnaryServerInterceptor(validator, opts...)
}

func ValidateStreamInterceptor(opts ...protovalidate_middleware.Option) grpc.StreamServerInterceptor {
	return protovalidate_middleware.StreamServerInterceptor(validator, opts...)
}
