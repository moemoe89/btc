package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	iDI "github.com/moemoe89/btc/internal/di"
	"github.com/moemoe89/btc/pkg/di"
)

func main() {
	server := iDI.GetBTCGRPCServer()

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
	di.CloseAll()

	log.Println("finished graceful shut down")
}
