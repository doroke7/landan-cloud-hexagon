package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、mysql.AbstractModel
// 的做法一致。Database 由呼叫端決定要選哪個 mongo database（bootstrap.NewMongo()
// 只回傳 *mongo.Client，選 database 這步留給組裝的地方做）。
type AbstractModel struct {
	Database *mongo.Database
	Context  context.Context
}

func NewAbstractModel(oContext context.Context, oDatabase *mongo.Database) *AbstractModel {
	return &AbstractModel{
		Database: oDatabase,
		Context:  oContext,
	}
}

// oNotDeletedAt 對應 mysql resource.sql 的軟刪除慣例：deleted_at 預設值不是
// NULL，而是這個「未來很久」的 sentinel（int32 最大 unix 時間），「還沒被刪除」
// 就用這個值表示，跟 mysql 那邊 Where("deleted_at = ?", "2038-01-19 03:14:07")
// 的做法一致，只是這裡存的是原生 bson 時間，不是字串。
var oNotDeletedAt = time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC)
