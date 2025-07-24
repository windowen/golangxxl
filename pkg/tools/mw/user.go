package mw

import (
	"context"

	"google.golang.org/grpc"

	"queueJob/pkg/constant"
)

func AddUserType() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// zlogger.Infow("add user type", zap.String("method", method))
		if arr, _ := ctx.Value(constant.RpcOpUserType).([]string); len(arr) > 0 {
			// zlogger.Infow("add user type", zap.String("method", method), zap.Strings("userType", arr))
			headers, _ := ctx.Value(constant.RpcCustomHeader).([]string)
			ctx = context.WithValue(ctx, constant.RpcCustomHeader, append(headers, constant.RpcOpUserType))
			ctx = context.WithValue(ctx, constant.RpcOpUserType, arr)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}
