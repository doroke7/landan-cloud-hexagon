package model

import (
	"context"

	helper "example/internal/helper"
)

// AbstractRepository 放 cache 這個 output adapter 共用的依賴——CacheHelper，
// 跟 mysql.AbstractRepository 持有 *gorm.DB 是同一種角色。
//
// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop 的做法一致：
// 程序收到中斷/終止訊號時，掛在這個 repository 底下的操作能跟著一起停止。
type AbstractRepository struct {
	CacheHelper *helper.CacheHelper
	Context     context.Context
}

func NewAbstractRepository(oContext context.Context, oCacheHelper *helper.CacheHelper) *AbstractRepository {
	return &AbstractRepository{
		CacheHelper: oCacheHelper,
		Context:     oContext,
	}
}
