package mw

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"queueJob/pkg/constant"
	"queueJob/pkg/protobuf/errinfo"
	"queueJob/pkg/tools/cast"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/zlogger"
)

func GrpcClient() grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(RpcClientInterceptor)
}

func RpcClientInterceptor(
	ctx context.Context,
	method string,
	req, resp interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) (err error) {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			zlogger.Errorf(tmpStr)
		}
	}()
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}

	err = invoker(ctx, method, req, resp, cc, opts...)
	if err == nil {
		// zlogger.Debug(ctx, "rpc client resp", "funcName", method, "resp", rpcString(resp))
		return nil
	}
	rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
	if !ok {
		return errs.ErrInternalServer.Wrap(err.Error())
	}
	sta := rpcErr.GRPCStatus()
	if sta.Code() == 0 {
		return errs.NewCodeError(errs.ServerInternalError, err.Error()).Wrap()
	}
	if details := sta.Details(); len(details) > 0 {
		errInfo, ok := details[0].(*errinfo.ErrorInfo)
		if ok {
			s := strings.Join(errInfo.Warp, "->") + errInfo.Cause
			return errs.NewCodeError(int(sta.Code()), sta.Message()).WithDetail(s).Wrap()
		}
	}
	return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
}

func getRpcContext(ctx context.Context, method string) (context.Context, error) {
	md := metadata.Pairs()
	if keys, _ := ctx.Value(constant.RpcCustomHeader).([]string); len(keys) > 0 {
		for _, key := range keys {
			val, ok := ctx.Value(key).([]string)
			if !ok {
				return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("ctx missing key %s", key))
			}
			if len(val) == 0 {
				return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("ctx key %s value is empty", key))
			}
			md.Set(key, val...)
		}
		md.Set(constant.RpcCustomHeader, keys...)
	}
	operationId, ok := ctx.Value(constant.OperationId).(string)
	if !ok {
		// zlogger.Warn(ctx, "ctx missing operationId", errors.New("ctx missing operationId"), "funcName", method)
		return nil, errs.ErrArgs.Wrap("ctx missing operationId")
	}
	md.Set(constant.OperationId, operationId)
	var checkArgs []string
	checkArgs = append(checkArgs, constant.OperationId, operationId)
	opUserId, ok := ctx.Value(constant.OpUserId).(int)
	if ok {
		md.Set(constant.OpUserId, cast.ToString(opUserId))
		checkArgs = append(checkArgs, constant.OpUserId, cast.ToString(opUserId))
	}
	opUserIdPlatformId, ok := ctx.Value(constant.OpUserPlatform).(string)
	if ok {
		md.Set(constant.OpUserPlatform, opUserIdPlatformId)
	}
	connId, ok := ctx.Value(constant.ConnId).(string)
	if ok {
		md.Set(constant.ConnId, connId)
	}
	countryCode, ok := ctx.Value(constant.CountryCode).(string)
	if ok {
		md.Set(constant.CountryCode, countryCode)
	}
	opTourist, ok := ctx.Value(constant.OpTourist).(int)
	if ok {
		md.Set(constant.OpTourist, cast.ToString(opTourist))
	}
	language, ok := ctx.Value(constant.Language).(string)
	if ok {
		md.Set(constant.Language, language)
	}

	location, ok := ctx.Value(constant.Location).(string)
	if ok {
		md.Set(constant.Location, location)
	}

	deviceId, ok := ctx.Value(constant.DeviceId).(string)
	if ok {
		md.Set(constant.DeviceId, deviceId)
	}

	return metadata.NewOutgoingContext(ctx, md), nil
}
