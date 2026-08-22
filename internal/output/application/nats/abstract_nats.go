package nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Conn 是從 bootstrap 注入的共用 nats 連線，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractNats struct {
	Context context.Context
	Conn    *nats.Conn
}

func NewAbstractNats(oContext context.Context, oConn *nats.Conn) *AbstractNats {
	return &AbstractNats{
		Context: oContext,
		Conn:    oConn,
	}
}
