package main

import (
	"bookstore/internal/server/book/server"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	server := server.NewBookServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("received signal to shutdown")
		cancel()
	}()

	if err := server.Run(ctx); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
