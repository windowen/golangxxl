package internal

import (
	"os"
	"strings"

	discov2 "queueJob/pkg/gozero/discov"
	"queueJob/pkg/gozero/netx"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

// NewRpcPubServer returns a Server.
func NewRpcPubServer(etcd discov2.EtcdConf, listenOn string,
	opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubListenOn := figureOutListenOn(listenOn)
		var pubOpts []discov2.PubOption
		if etcd.HasAccount() {
			pubOpts = append(pubOpts, discov2.WithPubEtcdAccount(etcd.User, etcd.Pass))
		}
		if etcd.HasTLS() {
			pubOpts = append(pubOpts, discov2.WithPubEtcdTLS(etcd.CertFile, etcd.CertKeyFile,
				etcd.CACertFile, etcd.InsecureSkipVerify))
		}
		if etcd.HasID() {
			pubOpts = append(pubOpts, discov2.WithId(etcd.ID))
		}
		pubClient := discov2.NewPublisher(etcd.Hosts, etcd.Key, pubListenOn, pubOpts...)
		return pubClient.KeepAlive()
	}
	server := keepAliveJob{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}

type keepAliveJob struct {
	registerEtcd func() error
	Server
}

func (s keepAliveJob) Start(fn RegisterFn) error {
	if err := s.registerEtcd(); err != nil {
		return err
	}

	return s.Server.Start(fn)
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
