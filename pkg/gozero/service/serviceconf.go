package service

import (
	"queueJob/pkg/gozero/load"
	prometheus2 "queueJob/pkg/gozero/prometheus"
	stat2 "queueJob/pkg/gozero/stat"
	trace2 "queueJob/pkg/gozero/trace"

	"queueJob/pkg/gozero/proc"
)

const (
	// DevMode means development mode.
	DevMode = "dev"
	// TestMode means test mode.
	TestMode = "test"
	// RtMode means regression test mode.
	RtMode = "rt"
	// PreMode means pre-release mode.
	PreMode = "pre"
	// ProMode means production mode.
	ProMode = "pro"
)

type (
	// DevServerConfig is type alias for devserver.Config

	// A ServiceConf is a service config.
	ServiceConf struct {
		Name       string
		Mode       string `json:",default=pro,options=dev|test|rt|pre|pro"`
		MetricsUrl string `json:",optional"`
		// Deprecated: please use DevServer
		Prometheus prometheus2.Config `json:",optional"`
		Telemetry  trace2.Config      `json:",optional"`
	}
)

// MustSetUp sets up the service, exits on error.
func (sc ServiceConf) MustSetUp() {
	err := sc.SetUp()
	if err != nil {
		panic(err)
	}
}

// SetUp sets up the service.
func (sc ServiceConf) SetUp() error {
	sc.initMode()
	prometheus2.StartAgent(sc.Prometheus)

	if len(sc.Telemetry.Name) == 0 {
		sc.Telemetry.Name = sc.Name
	}
	trace2.StartAgent(sc.Telemetry)
	proc.AddShutdownListener(func() {
		trace2.StopAgent()
	})

	if len(sc.MetricsUrl) > 0 {
		stat2.SetReportWriter(stat2.NewRemoteWriter(sc.MetricsUrl))
	}

	return nil
}

func (sc ServiceConf) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, RtMode, PreMode:
		load.Disable()
	}
}
