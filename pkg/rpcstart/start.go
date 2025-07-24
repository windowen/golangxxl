package rpcstart

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"liveJob/pkg/gozero/zrpc"

	"liveJob/pkg/common/config"
	"liveJob/pkg/gozero/discov"
	"liveJob/pkg/tools/errs"
	"liveJob/pkg/tools/mw"
	"liveJob/pkg/tools/network"
	"liveJob/pkg/zlogger"
)

func Start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(server *grpc.Server)) error {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, rpcRegisterName))
	zlogger.Infow("start", zap.String("register name", rpcKey), zap.Int("server port", rpcPort), zap.Int("prometheusPort:", prometheusPort))

	localIp, err := network.GetLocalIP()
	if err != nil {
		return errs.Wrap(err)
	}

	server := zrpc.MustNewServer(zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("%s:%d", localIp, rpcPort),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, rpcFn)

	server.AddOptions(mw.GrpcServer())

	server.Start()

	zlogger.Infof("%s RPC service shutdown", rpcRegisterName)
	return nil
}
