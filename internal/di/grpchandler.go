package di

import (
	"github.com/moemoe89/btc/internal/adapters/grpchandler"
)

// GetBTCGRPCHandler returns BTCServiceServer handler.
func GetBTCGRPCHandler() grpchandler.BTCServiceServer {
	return grpchandler.NewBTCHandler()
}
