package xetcd

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdConf struct {
	Hosts              []string
	User               string `yaml:",optional"`
	Pass               string `yaml:",optional"`
	CertFile           string `yaml:",optional"`
	CertKeyFile        string `yaml:",optional=CertFile"`
	CACertFile         string `yaml:",optional=CertFile"`
	InsecureSkipVerify bool   `yaml:",optional"`
}

// HasAccount returns if account provided.
func (c EtcdConf) HasAccount() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}

// HasTLS returns if TLS CertFile/CertKeyFile/CACertFile are provided.
func (c EtcdConf) HasTLS() bool {
	return len(c.CertFile) > 0 && len(c.CertKeyFile) > 0 && len(c.CACertFile) > 0
}

// DialClient dials an etcd cluster with given endpoints.
func DialClient(conf EtcdConf) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:            conf.Hosts,
		AutoSyncInterval:     time.Minute,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    5 * time.Second,
		DialKeepAliveTimeout: 5 * time.Second,
		RejectOldCluster:     true,
		PermitWithoutStream:  true,
	}
	if conf.HasAccount() {
		cfg.Username = conf.User
		cfg.Password = conf.Pass
	}
	if conf.HasTLS() {
		t, err := getTLS(conf.CertFile, conf.CertKeyFile, conf.CACertFile, conf.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}

		cfg.TLS = t
	}

	return clientv3.New(cfg)
}

func getTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
	if err != nil {
		return nil, err
	}

	caData, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: insecureSkipVerify,
	}, nil
}
