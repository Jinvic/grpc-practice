package interceptor

import (
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

func ValidateInterceptor(opts ...protovalidate_middleware.Option) grpc.UnaryServerInterceptor {
	return protovalidate_middleware.UnaryServerInterceptor(Validator, opts...)
}
