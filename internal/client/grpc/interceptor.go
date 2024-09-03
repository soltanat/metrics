package grpc

import (
	"context"
	"github.com/soltanat/metrics/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func LoggingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	latency := time.Since(start)

	l := logger.Get()

	var statusCode string
	if s, ok := status.FromError(err); ok {
		statusCode = s.Code().String()
	}

	l.Info().
		Str("method", method).
		Dur("latency", latency).
		Str("status", statusCode).
		Msg("Request")
	return err
}
