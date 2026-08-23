package bootstrap

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcd() (*clientv3.Client, error) {
	aEndpoints := make([]string, len(CONFIG.ETCD.HOSTS))
	for i, sHost := range CONFIG.ETCD.HOSTS {
		aEndpoints[i] = fmt.Sprintf("%s:%s", sHost, CONFIG.ETCD.PORTS[i])
	}

	oClient, oErr := clientv3.New(clientv3.Config{
		Endpoints:   aEndpoints,
		DialTimeout: time.Duration(CONFIG.ETCD.TIMEOUT) * time.Millisecond,
		Username:    CONFIG.ETCD.USER,
		Password:    CONFIG.ETCD.PASS,
	})

	log.Info("[INFO] ETCD 連線完成. ", strings.Join(aEndpoints, ","))

	return oClient, oErr
}
