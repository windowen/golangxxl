package clientinterceptors

import (
	"context"
	"path"
	"sync"
	"time"

	"google.golang.org/grpc"

	"queueJob/pkg/gozero/timex"

	"queueJob/pkg/gozero/lang"
	"queueJob/pkg/gozero/syncx"
	"queueJob/pkg/zlogger"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	notLoggingContentMethods sync.Map
	slowThreshold            = syncx.ForAtomicDuration(defaultSlowThreshold)
)

// DurationInterceptor is an interceptor that logs the processing time.
func DurationInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := timex.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		_, ok := notLoggingContentMethods.Load(method)
		if ok {
			zlogger.Errorf("fail - %s - %s", serverName, err.Error())
		} else {
			zlogger.Errorf("fail - %s - %v - %s", serverName, req, err.Error())
		}
	} else {
		elapsed := timex.Since(start)
		if elapsed > slowThreshold.Load() {
			_, ok := notLoggingContentMethods.Load(method)
			if ok {
				zlogger.Infof("[RPC] ok - slowcall - %s", serverName)
			} else {
				zlogger.Infof("[RPC] ok - slowcall - %s - %v - %v", serverName, req, reply)
			}
		}
	}

	return err
}

// DontLogContentForMethod disable logging content for given method.
func DontLogContentForMethod(method string) {
	notLoggingContentMethods.Store(method, lang.Placeholder)
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}
