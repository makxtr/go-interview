package common

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LoggingInterceptor is a server interceptor for logging
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Request received: %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("Request completed: %s in %v", info.FullMethod, time.Since(start))
	return resp, err
}

// AuthInterceptor is a client interceptor for authentication
func AuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add auth token to metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")
	return invoker(ctx, method, req, reply, cc, opts...)
}
