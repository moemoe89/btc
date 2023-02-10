package grpchandler

import (
	"context"

	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// TODO: Move both Check and Watch method as a pkg.

// Check checks if the gRPC handler is ready to receive a request from client.
// https://github.com/americanexpress/grpc-k8s-health-check/blob/master/server-grpc/server.go#L57
func (h *btcHandler) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch serves status of the requested services whenever the status changes.
func (h *btcHandler) Watch(*health.HealthCheckRequest, health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
