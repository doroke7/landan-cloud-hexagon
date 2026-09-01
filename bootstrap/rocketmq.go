package bootstrap

import (
	"strings"

	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/charmbracelet/log"
)

// RocketMQ 用 NameServer 定址，跟其他 broker 常見的 host:port 連線方式不同；
// 這裡回傳的是已經 Start() 好的 Producer，讓 output adapter 可以直接送訊息。
func NewRocketmq() (rocketmq.Producer, error) {
	aOptions := []producer.Option{
		producer.WithNameServer(primitive.NamesrvAddr(CONFIG.ROCKETMQ.NAME_SERVERS)),
		producer.WithGroupName(CONFIG.ROCKETMQ.GROUP),
		producer.WithRetry(2),
	}

	if CONFIG.ROCKETMQ.ACCESS_KEY != "" {
		aOptions = append(aOptions, producer.WithCredentials(primitive.Credentials{
			AccessKey: CONFIG.ROCKETMQ.ACCESS_KEY,
			SecretKey: CONFIG.ROCKETMQ.SECRET_KEY,
		}))
	}

	oProducer, oErr := rocketmq.NewProducer(aOptions...)
	if oErr != nil {
		return nil, oErr
	}

	if oErr := oProducer.Start(); oErr != nil {
		return nil, oErr
	}

	log.Info("[INFO] ROCKETMQ 連線完成.", "addr", strings.Join(CONFIG.ROCKETMQ.NAME_SERVERS, ","))

	return oProducer, nil
}
