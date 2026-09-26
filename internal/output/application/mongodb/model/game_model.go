package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
)

// oDeletedAtZero 跟 mysql/resource adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、
// name+deleted_at、deleted_at、game_type_id+deleted_at 供查詢用。
func NewGameModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) (outputPortAnyModel.GameModel, error) {
	oCollection := oAbstractMongodb.Database.Collection("games")
	oIndexView := oCollection.Indexes()

	oKeyIndexOptions := options.Index()
	oKeyIndexOptions.SetUnique(true)
	oKeyIndexOptions.SetName("games-k")

	oNameDeletedAtIndexOptions := options.Index()
	oNameDeletedAtIndexOptions.SetName("games-n-da")

	oDeletedAtIndexOptions := options.Index()
	oDeletedAtIndexOptions.SetName("games-da")

	oGameTypeIdDeletedAtIndexOptions := options.Index()
	oGameTypeIdDeletedAtIndexOptions.SetName("games-gti-da")

	aIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: oKeyIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: oNameDeletedAtIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: oDeletedAtIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "game_type_id", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: oGameTypeIdDeletedAtIndexOptions,
		},
	}

	if _, oErr := oIndexView.CreateMany(oAbstractMongodb.Context, aIndexModels); oErr != nil {
		return nil, oErr
	}

	oCounters := oAbstractMongodb.Database.Collection("counters")

	return &GameModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oCounters,
	}, nil
}

func (oSelf *GameModel) nextId() (uint, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "game"},
		bson.M{"$inc": bson.M{"seq": 1}},
		oUpdateOptions,
	)

	var oCounter struct {
		Seq uint `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *GameModel) AddOne(oGame *domain.GameVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.Game{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oGame.GameTypeId != nil {
		oNew.GameTypeId = *oGame.GameTypeId
	}
	if oGame.Key != nil {
		oNew.Key = *oGame.Key
	}
	if oGame.Name != nil {
		oNew.Name = *oGame.Name
	}
	if oGame.Description != nil {
		oNew.Description = *oGame.Description
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {
	var oGame domain.Game

	oResult := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"key":        sKey,
		"deleted_at": oDeletedAtZero,
	})
	oErr := oResult.Decode(&oGame)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oGame.GameTypeId != nil {
		oSet["game_type_id"] = *oGame.GameTypeId
	}
	if oGame.Key != nil {
		oSet["key"] = *oGame.Key
	}
	if oGame.Name != nil {
		oSet["name"] = *oGame.Name
	}
	if oGame.Description != nil {
		oSet["description"] = *oGame.Description
	}

	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId},
		bson.M{"$set": oSet},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.MatchedCount == 0 {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *GameModel) RemoveOneById(iId uint64) error {
	oNow := time.Now()
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": oNow}},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.ModifiedCount == 0 {
		return errors.New("0 rows deleted")
	}

	return nil
}

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint64) (uint, error) {
	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, bson.M{
		"game_type_id": iGameTypeId,
		"deleted_at":   oDeletedAtZero,
	})
	if oErr != nil {
		return 0, oErr
	}

	return uint(iCount), nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
