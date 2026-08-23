package mosquitto

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Client 是從 bootstrap 注入的共用 MQTT client，
// 跟 mysql.AbstractModel 持有 *gorm.DB 是同一種角色。
type AbstractMosquitto struct {
	Context context.Context
	Client  mqtt.Client
}

func NewAbstractMosquitto(oContext context.Context, oClient mqtt.Client) *AbstractMosquitto {
	return &AbstractMosquitto{
		Context: oContext,
		Client:  oClient,
	}
}
