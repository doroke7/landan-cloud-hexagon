package etcd

import (
	"context"

	helper "example/internal/helper"
)

// AbstractRepository 放 etcd 這個 output adapter 共用的依賴——EtcdHelper，
// 跟 cache.AbstractCache 持有 *helper.CacheHelper 是同一種角色。
//
// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop 的做法一致：
// 程序收到中斷/終止訊號時，掛在這個 repository 底下的操作能跟著一起停止。
type AbstractEtcd struct {
	EtcdHelper *helper.EtcdHelper
	Context    context.Context
}

func NewAbstractEtcd(oContext context.Context, oEtcdHelper *helper.EtcdHelper) *AbstractEtcd {
	return &AbstractEtcd{
		EtcdHelper: oEtcdHelper,
		Context:    oContext,
	}
}
