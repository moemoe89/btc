package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moemoe89/btc/internal/di"
)

func main() {
	server := di.GetBTCGRPCServer()

	log.Println("BTC service is ready")

	go func() {
		// Run() keeps its process until receiving any error
		if err := server.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %d received, shutting down gracefully...", <-quit)
	server.GracefulStop()

	log.Println("finished graceful shut down")
}
