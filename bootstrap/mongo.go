package bootstrap

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongo() (*mongo.Client, error) {
	sURI := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		CONFIG.MONGODB.PROTOCOL,
		CONFIG.MONGODB.USER,
		CONFIG.MONGODB.PASSWORD,
		CONFIG.MONGODB.HOST,
		CONFIG.MONGODB.PORT,
		CONFIG.MONGODB.NAME,
	)

	oOptions := options.Client().
		ApplyURI(sURI).
		SetMaxPoolSize(CONFIG.MONGODB.MAX_POOL_SIZE).
		SetMinPoolSize(CONFIG.MONGODB.MIN_POOL_SIZE)

	oMongoConnectionPool, oErr := mongo.Connect(oOptions)

	return oMongoConnectionPool, oErr
}
