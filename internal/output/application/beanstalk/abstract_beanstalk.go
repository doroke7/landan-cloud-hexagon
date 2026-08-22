package beanstalk

import (
	"context"

	"github.com/beanstalkd/go-beanstalk"
	beanstalkd "github.com/beanstalkd/go-beanstalk"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Conn 是從 bootstrap 注入的共用 beanstalkd 連線，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractBeanstalk struct {
	Context context.Context
	Conn    *beanstalk.Conn
}

func NewAbstractBeanstalk(oContext context.Context, oConn *beanstalk.Conn) *AbstractBeanstalk {
	return &AbstractBeanstalk{
		Context: oContext,
		Conn:    oConn,
	}
}

func (oSelf *AbstractBeanstalk) Tube(sName string) *beanstalk.Tube {
	return beanstalkd.NewTube(oSelf.Conn, sName)
}
