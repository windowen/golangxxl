package metric

import (
	"liveJob/pkg/gozero/prometheus"
)

// A VectorOpts is a general configuration.
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

func update(fn func()) {
	if !prometheus.Enabled() {
		return
	}

	fn()
}
