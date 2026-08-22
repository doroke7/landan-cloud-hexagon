package zeromq

import (
	"context"

	"github.com/go-zeromq/zmq4"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Socket 是從 bootstrap 注入的共用 PUB socket，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractZeromq struct {
	Context context.Context
	Socket  zmq4.Socket
}

func NewAbstractZeromq(oContext context.Context, oSocket zmq4.Socket) *AbstractZeromq {
	return &AbstractZeromq{
		Context: oContext,
		Socket:  oSocket,
	}
}
