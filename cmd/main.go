package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/moemoe89/btc/internal/di"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := di.GetBTCGRPCServer()

	log.Println("BTC service is ready")

	go func() {
		// Run() keeps its process until receiving any error
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %d received, shutting down gracefully...", <-quit)
	server.GracefulStop()

	log.Println("finished graceful shut down")
}
