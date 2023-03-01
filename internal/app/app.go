package app

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var EtcdClientCli *clientv3.Client

func InitApplication() (err error, cleanup func()) {

	EtcdClientCli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}

	return nil, func() {
		_ = EtcdClientCli.Close()

		log.Println("Application cleanup")
	}
}
