package apiresp

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/zlogger"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client C,
	c *gin.Context,
) {
	var req A
	if err := c.BindJSON(&req); err != nil {
		zlogger.Warnw("gin bind json error", zap.Error(err), zap.Any("req", req))
		GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap()) // 参数错误
		return
	}
	if err := Validate(&req); err != nil {
		GinError(c, err) // 参数校验失败
		return
	}
	data, err := rpc(client, c, &req)
	if err != nil {
		GinError(c, err) // RPC调用失败
		return
	}

	GinSuccess(c, data) // 成功
}
