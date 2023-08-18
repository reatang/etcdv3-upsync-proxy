package app

import "github.com/reatang/etcdv3-upsync-proxy/pkg/xetcd"

type ServerConf struct {
	ListenOn string
	Etcd     xetcd.EtcdConf
}
