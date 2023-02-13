package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	iDI "github.com/moemoe89/btc/internal/di"
	"github.com/moemoe89/btc/pkg/di"

	"go.uber.org/zap"
)

func main() {
	logger := iDI.GetLogger()

	server := iDI.GetBTCGRPCServer()

	logger.Info("BTC service is ready")

	go func() {
		// Run() keeps its process until receiving any error
		if err := server.Run(); err != nil {
			logger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	logger.Info(fmt.Sprintf("SIGNAL %d received, shutting down gracefully...", <-quit))
	di.CloseAll()

	logger.Info("finished graceful shut down")
}
