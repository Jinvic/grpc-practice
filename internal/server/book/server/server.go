package server

import (
	bookv1 "bookstore/api/book/v1"
	"bookstore/config"
	"bookstore/internal/server/book/service"
	"bookstore/internal/server/common/interceptor"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/samber/do/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type BookServer struct {
}

func NewBookServer() *BookServer {
	return &BookServer{}
}

func (s *BookServer) Run(ctx context.Context) error {
	injector := do.New()
	if err := InitInjector(injector); err != nil {
		return fmt.Errorf("failed to initialize injector: %w", err)
	}

	port := do.MustInvoke[*config.Config](injector).Services.Book.Port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer lis.Close()

	if err := interceptor.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize interceptor: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.ValidateUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptor.ValidateStreamInterceptor(),
		),
	)
	reflection.Register(grpcServer)

	bookService := do.MustInvoke[*service.BookService](injector)
	bookv1.RegisterBookServiceServer(grpcServer, bookService)

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("server listening at port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Println("server stopped")
			return nil
		}
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
