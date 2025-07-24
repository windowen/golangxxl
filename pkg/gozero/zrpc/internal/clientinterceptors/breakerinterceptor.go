package clientinterceptors

import (
	"context"
	"path"

	"google.golang.org/grpc"

	"liveJob/pkg/gozero/zrpc/internal/codes"

	"liveJob/pkg/gozero/breaker"
)

// BreakerInterceptor is an interceptor that acts as a circuit breaker.
func BreakerInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptableCtx(ctx, breakerName, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, codes.Acceptable)
}
