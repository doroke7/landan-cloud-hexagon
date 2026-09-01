package bootstrap

import (
	"fmt"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongo() (*mongo.Client, error) {
	// authSource 跟實際要用的資料庫（NAME）分開設定：連線字串路徑那段預設同時是
	// 「預設資料庫」跟「預設 authSource」，但帳號通常是在另一個資料庫（例如 admin
	// 或建帳號當下的那個庫）建立的，兩者混用會導致 SCRAM 認證失敗。

	sURI := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?authSource=%s",
		CONFIG.MONGODB.PROTOCOL,
		CONFIG.MONGODB.USER,
		CONFIG.MONGODB.PASSWORD,
		CONFIG.MONGODB.HOST,
		CONFIG.MONGODB.PORT,
		CONFIG.MONGODB.NAME,
		CONFIG.MONGODB.AUTH_SOURCE,
	)

	oOptions := options.Client().
		ApplyURI(sURI).
		SetMaxPoolSize(CONFIG.MONGODB.MAX_POOL_SIZE).
		SetMinPoolSize(CONFIG.MONGODB.MIN_POOL_SIZE)

	oMongoConnectionPool, oErr := mongo.Connect(oOptions)

	log.Info("[INFO] MONGODB 連線完成.", "addr", CONFIG.MONGODB.HOST+":"+CONFIG.MONGODB.PORT)

	return oMongoConnectionPool, oErr
}
