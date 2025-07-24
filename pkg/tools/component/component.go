package component

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"queueJob/pkg/common/config"
	"queueJob/pkg/zlogger"
)

var (
	MaxConnectTimes = 100
)

var ConnectClient = &DiscoveryClient{}

type DiscoveryClient struct {
	etcdConn *clientv3.Client
}

func ComponentCheck(cfgPath string, discoveryName string, hide bool) error {
	if discoveryName != "k8s" {
		if _, err := checkNewDiscoveryClient(discoveryName, hide); err != nil {
			errorPrint(fmt.Sprintf("%v.Please check if your  server has started", err.Error()), hide)
			return err
		}
	}

	return nil
}

func errorPrint(s string, hide bool) {
	if !hide {
		fmt.Printf("\x1b[%dm%v\x1b[0m\n", 31, s)
	}
}

func successPrint(s string, hide bool) {
	if !hide {
		fmt.Printf("\x1b[%dm%v\x1b[0m\n", 32, s)
	}
}

func newDiscoveryClient(discoveryName string) (*DiscoveryClient, error) {
	if discoveryName == "etcd" {
		cfg := clientv3.Config{
			Endpoints:   config.Config.Etcd.Addr, // Etcd 服务器地址
			DialTimeout: 5 * time.Second,         // 连接超时设置
		}
		con, err := clientv3.New(cfg)
		if err != nil {
			panic(err)
		}
		ConnectClient.etcdConn = con

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = con.Get(ctx, "key")
		if err != nil {
			zlogger.Errorw("获取ETCD服务注册发现失败！", zap.Error(err))
			panic(err)
		}

		fmt.Println("etcd addr=", config.Config.Etcd.Addr)
	}
	return ConnectClient, nil
}

func checkNewDiscoveryClient(discoveryName string, hide bool) (*DiscoveryClient, error) {
	conn, err := newDiscoveryClient(discoveryName)
	if err != nil {
		if conn != nil {
			if discoveryName == "etcd" {
				conn.etcdConn.Close()
			}
		}
	}
	successPrint(fmt.Sprintf("%s starts successfully", discoveryName), hide)
	return conn, nil
}
