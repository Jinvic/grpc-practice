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

func ValidateInterceptor(opts ...protovalidate_middleware.Option) grpc.UnaryServerInterceptor {
	return protovalidate_middleware.UnaryServerInterceptor(validator, opts...)
}
