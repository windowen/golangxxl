package rpcstart

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"queueJob/pkg/gozero/zrpc"

	"queueJob/pkg/common/config"
	"queueJob/pkg/gozero/discov"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/tools/mw"
	"queueJob/pkg/tools/network"
	"queueJob/pkg/zlogger"
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
