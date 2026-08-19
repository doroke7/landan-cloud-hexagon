package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (redis.UniversalClient, error) {
	var oRedisConnectionPool redis.UniversalClient

	if CONFIG.REDIS.CLUSTER {
		aAddrs := make([]string, len(CONFIG.REDIS.HOSTS))
		for i, sHost := range CONFIG.REDIS.HOSTS {
			aAddrs[i] = fmt.Sprintf("%s:%s", sHost, CONFIG.REDIS.PORTS[i])
		}

		oRedisConnectionPool = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    aAddrs,
			Username: CONFIG.REDIS.USERNAME,
			Password: CONFIG.REDIS.PASSWORD,
		})
	} else {
		sAddr := fmt.Sprintf("%s:%s", CONFIG.REDIS.HOST, CONFIG.REDIS.PORT)

		oRedisConnectionPool = redis.NewClient(&redis.Options{
			Addr:     sAddr,
			Username: CONFIG.REDIS.USERNAME,
			Password: CONFIG.REDIS.PASSWORD,
			DB:       CONFIG.REDIS.DB,
		})
	}

	if err := oRedisConnectionPool.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return oRedisConnectionPool, nil
}
