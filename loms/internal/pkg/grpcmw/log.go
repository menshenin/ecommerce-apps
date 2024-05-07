// Package grpcmw Middleware для GRPC
package grpcmw

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

// Logged Логирующая Middleware
func Logged(logger *slog.Logger) func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.InfoContext(ctx, "REQUEST:", "data", req)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.ErrorContext(ctx, "ERROR:", "data", err)
			return resp, err
		}
		logger.InfoContext(ctx, "RESPONSE:", "data", resp)
		return resp, err
	}
}
