package middleware

import (
	internal "liveJob/pkg/context"
	"liveJob/pkg/xxl"
)

// CustomMiddleware 自定义中间件
func CustomMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(ctx *internal.Context, param *xxl.RunReq) string {
		// startTime := time.Now()
		// ctx.Infof("[middleware] Start at: %v", startTime)
		// ctx.Console("<<<< 执行脚本: %s >>>> <br>", param.ExecutorHandler)
		// ctx.Console(
		// 	"<<<< 脚本参数 >>>> <br> %s <br>------------------------------------------------------------------------------------------------<br>",
		// 	param.ExecutorParams)
		res := tf(ctx, param)
		// since := time.Since(startTime)
		// ctx.Infof("[middleware] End... execution time: %v", since)
		return res
	}
}
