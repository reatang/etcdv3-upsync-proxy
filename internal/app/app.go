package app

import (
	"log"

	"github.com/reatang/etcdv3-upsync-proxy/pkg/xetcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdClient *clientv3.Client

func InitApplication(c ServerConf) (err error, cleanup func()) {
	EtcdClient, err = xetcd.DialClient(c.Etcd)
	if err != nil {
		return
	}

	return nil, func() {
		_ = EtcdClient.Close()

		log.Println("Application cleanup")
	}
}
