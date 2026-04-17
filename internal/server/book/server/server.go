package server

import (
	bookv1 "bookstore/api/book/v1"
	"bookstore/internal/common/interceptor"
	"context"
	"fmt"
	"log"
	"net"

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
			interceptor.ValidateInterceptor(),
		),
	)
	reflection.Register(grpcServer)
	bookv1.RegisterBookServiceServer(grpcServer, BuildBookService())

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
