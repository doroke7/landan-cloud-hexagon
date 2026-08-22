package bootstrap

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcd() (*clientv3.Client, error) {
	aEndpoints := make([]string, len(CONFIG.ETCD.HOSTS))
	for i, sHost := range CONFIG.ETCD.HOSTS {
		aEndpoints[i] = fmt.Sprintf("%s:%s", sHost, CONFIG.ETCD.PORTS[i])
	}

	return clientv3.New(clientv3.Config{
		Endpoints:   aEndpoints,
		DialTimeout: time.Duration(CONFIG.ETCD.TIMEOUT) * time.Millisecond,
		Username:    CONFIG.ETCD.USER,
		Password:    CONFIG.ETCD.PASS,
	})
}
