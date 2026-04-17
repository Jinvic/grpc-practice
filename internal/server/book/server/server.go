package server

import (
	bookv1 "bookstore/api/book/v1"
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
	port int
}

func NewBookServer(port int) *BookServer {
	return &BookServer{port: port}
}

func (s *BookServer) Run(ctx context.Context) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
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

	if err := InitInjector(); err != nil {
		return fmt.Errorf("failed to initialize injector: %w", err)
	}
	bookService := do.MustInvoke[*service.BookService](injector)
	bookv1.RegisterBookServiceServer(grpcServer, bookService)

	go func() {
		<-ctx.Done()
		log.Println("shutting down server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("server listening at port %d", s.port)
	if err := grpcServer.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Println("server stopped")
			return nil
		}
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
