package di

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/moemoe89/btc/pkg/di"
	"github.com/moemoe89/btc/pkg/trace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type wrapProvider struct {
	provider trace.Provider
}

func (w *wrapProvider) Close() error {
	return w.provider.Close(context.Background())
}

func GetTrace() []grpc.ServerOption {
	// Bootstrap tracer.
	prv, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: os.Getenv("OTEL_AGENT"),
		ServiceName:    "BTC Server",
		ServiceVersion: "1.0.0",
		Environment:    os.Getenv("APP_ENV"),
		Disabled:       false,
	})
	if err != nil {
		log.Fatal("trace new provider", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	}

	var c io.Closer = &wrapProvider{
		provider: prv,
	}

	di.RegisterCloser("Trace Provider", c)

	return opts
}
