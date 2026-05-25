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
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if s.cfg.Services.Book.EnableHTTP {
		go func() {
			time.Sleep(1 * time.Second) // wait for grpc server to start
			if err := s.runHTTPGateway(ctx); err != nil {
				log.Printf("failed to run http gateway: %v", err)
			}
		}()
	}

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
		log.Println("shutting down grpc server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("grpc server listening at port %d", s.cfg.Services.Book.Port)
	if err := grpcServer.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			log.Println("grpc server stopped")
			return nil
		}
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *BookServer) runHTTPGateway(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	grpcEndpoint := fmt.Sprintf("localhost:%d", s.cfg.Services.Book.Port)
	err := bookv1.RegisterBookServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	httpAddr := fmt.Sprintf(":%d", s.cfg.Services.Book.HTTPPort)
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("shutting down http gateway...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("failed to shutdown http gateway: %v", err)
		}
	}()

	log.Printf("http gateway listening at port %d", s.cfg.Services.Book.HTTPPort)
	if err := httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("http gateway closed")
			return nil
		}
		return fmt.Errorf("failed to serve http gateway: %w", err)
	}
	return nil
}
