package di

import (
	"log"
	"os"

	rpc "github.com/moemoe89/btc/api/go/grpc"
	"github.com/moemoe89/btc/pkg/di"
	"github.com/moemoe89/btc/pkg/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GetBTCGRPCServer returns gRPC server instance for BTC service.
func GetBTCGRPCServer() server.Server {
	s, err := server.NewGRPCServer(os.Getenv("SERVER_PORT"), func(server *grpc.Server) {
		h := GetBTCGRPCHandler()
		rpc.RegisterBTCServiceServer(server, h)
		grpc_health_v1.RegisterHealthServer(server, h)
	})
	if err != nil {
		log.Fatal(err)
	}

	di.RegisterCloser("gRPC server", di.NewCloser(s.GracefulStop))

	return s
}
