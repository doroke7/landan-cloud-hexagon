package activemq

import (
	"context"

	"github.com/go-stomp/stomp/v3"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Conn 是從 bootstrap 注入的共用 activemq（STOMP）連線，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractActivemq struct {
	Context context.Context
	Conn    *stomp.Conn
}

func NewAbstractActivemq(oContext context.Context, oConn *stomp.Conn) *AbstractActivemq {
	return &AbstractActivemq{
		Context: oContext,
		Conn:    oConn,
	}
}
