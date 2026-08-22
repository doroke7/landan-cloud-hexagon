package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"

	bootstrap "example/bootstrap"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractModel 做法一致。
// Database 是從 bootstrap 注入的共用 mongo 連線取出的資料庫控制代碼（CONFIG.MONGODB.NAME），
// 跟 mysql.AbstractMysql 持有 *gorm.DB 是同一種角色。
type AbstractMongodb struct {
	Context  context.Context
	Client   *mongo.Client
	Database *mongo.Database
}

func NewAbstractMongodb(oContext context.Context, oClient *mongo.Client) *AbstractMongodb {
	return &AbstractMongodb{
		Context:  oContext,
		Client:   oClient,
		Database: oClient.Database(bootstrap.CONFIG.MONGODB.NAME),
	}
}
