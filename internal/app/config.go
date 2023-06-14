package app

import "github.com/reatang/etcdv3-upsync-proxy/pkg/xetcd"

type ServerConf struct {
	ListenOn string         `yaml:"ListenOn"`
	Etcd     xetcd.EtcdConf `yaml:"Etcd"`
}
