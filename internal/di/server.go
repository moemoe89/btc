package di

import (
	rpc "github.com/moemoe89/btc/api/go/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GetBTCGRPCServer returns gRPC server instance for BTC service.
func GetBTCGRPCServer() *grpc.Server {
	s := grpc.NewServer()

	h := GetBTCGRPCHandler()
	rpc.RegisterBTCServiceServer(s, h)
	grpc_health_v1.RegisterHealthServer(s, h)

	return s
}
