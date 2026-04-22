package main

import (
	"bookstore/internal/pkg/config"
	"bookstore/internal/pkg/logger"
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

	injector := do.New()
	if err := server.InitInjector(injector); err != nil {
		log.Fatalf("failed to initialize injector: %v", err)
	}

	cfg := do.MustInvoke[*config.Config](injector)
	cleanup, err := logger.InitLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer cleanup()

	bookServer, err := server.NewBookServer(injector)
	if err != nil {
		log.Fatalf("failed to create book server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
