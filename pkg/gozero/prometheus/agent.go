package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"liveJob/pkg/gozero/syncx"
	"liveJob/pkg/gozero/threading"
	"liveJob/pkg/zlogger"
)

var (
	once    sync.Once
	enabled syncx.AtomicBool
)

// Enabled returns if prometheus is enabled.
func Enabled() bool {
	return enabled.True()
}

// Enable enables prometheus.
func Enable() {
	enabled.Set(true)
}

// StartAgent starts a prometheus agent.
func StartAgent(c Config) {
	if len(c.Host) == 0 {
		return
	}

	once.Do(func() {
		enabled.Set(true)
		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			zlogger.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				zlogger.Error(err)
			}
		})
	})
}
