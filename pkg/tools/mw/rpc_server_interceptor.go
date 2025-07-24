package mw

import (
	"context"
	"fmt"
	"math"
	"runtime"

	"queueJob/pkg/constant"
	"queueJob/pkg/zlogger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"queueJob/pkg/protobuf/errinfo"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/tools/mw/specialerror"
)

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

func RpcServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			zlogger.Errorf(tmpStr)
		}
	}()
	// funcName := info.FullMethod
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "missing metadata").Err()
	}
	if keys := md.Get(constant.RpcCustomHeader); len(keys) > 0 {
		for _, key := range keys {
			values := md.Get(key)
			if len(values) == 0 {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("missing metadata key %s", key)).Err()
			}
			ctx = context.WithValue(ctx, key, values)
		}
	}
	args := make([]string, 0, 4)
	if opts := md.Get(constant.OperationId); len(opts) != 1 || opts[0] == "" {
		return nil, status.New(codes.InvalidArgument, "operationID error").Err()
	} else {
		args = append(args, constant.OperationId, opts[0])
		ctx = context.WithValue(ctx, constant.OperationId, opts[0])
	}
	if opts := md.Get(constant.OpUserId); len(opts) == 1 {
		args = append(args, constant.OpUserId, opts[0])
		ctx = context.WithValue(ctx, constant.OpUserId, opts[0])
	}
	if opts := md.Get(constant.OpUserPlatform); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.OpUserPlatform, opts[0])
	}
	if opts := md.Get(constant.ConnId); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.ConnId, opts[0])
	}
	if opts := md.Get(constant.CountryCode); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.CountryCode, opts[0])
	}
	if opts := md.Get(constant.Language); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.Language, opts[0])
	}
	resp, err = func() (interface{}, error) {
		return handler(ctx, req)
	}()
	if err == nil {
		// zlogger.Debug(ctx, "rpc server resp", "funcName", funcName, "resp", rpcString(resp))
		return resp, nil
	}
	unwrap := errs.Unwrap(err)
	codeErr := specialerror.ErrCode(unwrap)
	if codeErr == nil {
		// zlogger.Error(ctx, "rpc InternalServer error", err, "req", req)
		codeErr = errs.ErrInternalServer
	}
	code := codeErr.Code()
	if code <= 0 || code > math.MaxUint32 {
		// zlogger.Error(ctx, "rpc UnknownError", err, "rpc UnknownCode:", code)
		code = errs.ServerInternalError
	}
	grpcStatus := status.New(codes.Code(code), codeErr.Msg())
	errInfo := &errinfo.ErrorInfo{Cause: err.Error()}
	details, err := grpcStatus.WithDetails(errInfo)
	if err != nil {
		// zlogger.Warn("rpc server resp WithDetails error", err, "funcName", funcName)
		return nil, errs.Wrap(err)
	}
	return nil, details.Err()
}

func GrpcServer() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(RpcServerInterceptor)
}
