package main

import (
	"bookstore/internal/pkg/config"
	"bookstore/internal/pkg/logger"
	"bookstore/internal/pkg/otel"
	"bookstore/internal/server/book/server"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/do/v2"
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	injector := do.New()
	if err := server.InitInjector(injector); err != nil {
		log.Fatalf("failed to initialize injector: %v", err)
	}

	cfg := do.MustInvoke[*config.Config](injector)
	logger.InitLogger(&cfg.Logging, cfg.Services.Book.LogFile)
	shutdownOtel, err := otel.InitOtel(ctx, &cfg.Otel, cfg.Services.Book.Name, cfg.Services.Book.Version, cfg.Environment)
	if err != nil {
		log.Fatalf("failed to initialize otel: %v", err)
	}
	defer func() {
		if err := shutdownOtel(context.Background()); err != nil {
			log.Fatalf("failed to shutdown otel: %v", err)
		}
	}()

	bookServer, err := server.NewBookServer(injector)
	if err != nil {
		log.Fatalf("failed to create book server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("received signal to shutdown")
		cancel()
	}()

	if err := bookServer.Run(ctx); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
