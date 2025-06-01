package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/yourusername/api/proto"
)

// RunGatewayServer starts the HTTP REST gateway
func RunGatewayServer(grpcAddr, httpAddr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}

	log.Printf("🌐 gRPC-Gateway HTTP server started on %s (proxy to %s)", httpAddr, grpcAddr)
	return http.ListenAndServe(httpAddr, mux)
}
