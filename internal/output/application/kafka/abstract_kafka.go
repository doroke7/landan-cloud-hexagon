package outputApplicationKafka

import (
	"context"

	"github.com/IBM/sarama"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Producer 是從 bootstrap 注入的共用 sarama Client 建出來的 SyncProducer，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractKafka struct {
	Context  context.Context
	Client   sarama.Client
	Producer sarama.SyncProducer
}

func NewAbstractKafka(oContext context.Context, oClient sarama.Client) (*AbstractKafka, error) {
	oProducer, err := sarama.NewSyncProducerFromClient(oClient)
	if err != nil {
		return nil, err
	}

	return &AbstractKafka{
		Context:  oContext,
		Client:   oClient,
		Producer: oProducer,
	}, nil
}
