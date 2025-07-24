package sysx

import (
	"os"

	"liveJob/pkg/gozero/stringx"
)

var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = stringx.RandId()
	}
}

// Hostname returns the name of the host, if no hostname, a random id is returned.
func Hostname() string {
	return hostname
}
