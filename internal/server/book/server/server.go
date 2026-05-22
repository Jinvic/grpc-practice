package server

import (
	bookv1 "bookstore/api/book/v1"
	"bookstore/internal/pkg/config"
	"bookstore/internal/server/book/service"
	"bookstore/internal/server/common/interceptor"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type BookServer struct {
	cfg         *config.Config
	bookService *service.BookService
}

func NewBookServer(i do.Injector) (*BookServer, error) {
	cfg := do.MustInvoke[*config.Config](i)
	bookService := do.MustInvoke[*service.BookService](i)
	return &BookServer{
		cfg:         cfg,
		bookService: bookService,
	}, nil
}

func (s *BookServer) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Services.Book.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	if err := interceptor.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize interceptor: %w", err)
	}

	loggingOptions := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.ValidateUnaryInterceptor(),
			interceptor.LoggingUnaryInterceptor(loggingOptions...),
		),
		grpc.ChainStreamInterceptor(
			interceptor.ValidateStreamInterceptor(),
			interceptor.LoggingStreamInterceptor(loggingOptions...),
		),
	)
	reflection.Register(grpcServer)

	bookv1.RegisterBookServiceServer(grpcServer, s.bookService)

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("server listening at port %d", s.cfg.Services.Book.Port)
	if err := grpcServer.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Println("server stopped")
			return nil
		}
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
