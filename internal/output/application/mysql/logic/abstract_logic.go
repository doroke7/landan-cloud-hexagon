package mysql

import (
	"context"

	pkg "example/pkg"

	"gorm.io/gorm"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractLogic struct {
	DB      *gorm.DB
	Context context.Context
	*pkg.Aop
}

func NewAbstractLogic(oContext context.Context, oDb *gorm.DB, oAop *pkg.Aop) *AbstractLogic {

	return &AbstractLogic{
		DB:      oDb,
		Context: oContext,
		Aop:     oAop,
	}
}
