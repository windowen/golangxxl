package service

import (
	"liveJob/pkg/gozero/service"
	"liveJob/pkg/zlogger"
)

var group *service.ServiceGroup

func init() {
	group = service.NewServiceGroup()
}

func RegisterService(service service.Service) {
	group.Add(service)
}

func Start(info string) {
	group.Start()
	zlogger.Infof("service start: %s", info)
}

func Stop(info string) {
	group.Stop()
	zlogger.Infof("service stop: %s", info)
}
