package outputApplicationPulsar

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Client 是從 bootstrap 注入的共用 pulsar 連線，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractPulsar struct {
	Context context.Context
	Client  pulsar.Client
}

func NewAbstractPulsar(oContext context.Context, oClient pulsar.Client) *AbstractPulsar {
	return &AbstractPulsar{
		Context: oContext,
		Client:  oClient,
	}
}
