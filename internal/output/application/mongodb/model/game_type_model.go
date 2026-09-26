package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	outputApplicationMongodb "example/internal/output/application/mongodb"
	pkgInput "example/pkg/input"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type GameTypeModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、name+deleted_at、deleted_at 供查詢用。
func NewGameTypeModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) (outputPortAnyModel.GameTypeModel, error) {
	oCollection := oAbstractMongodb.Database.Collection("game_types")
	oIndexView := oCollection.Indexes()

	oKeyIndexOptions := options.Index()
	oKeyIndexOptions.SetUnique(true)
	oKeyIndexOptions.SetName("game_types-k")

	oNameDeletedAtIndexOptions := options.Index()
	oNameDeletedAtIndexOptions.SetName("game_types-n-da")

	oDeletedAtIndexOptions := options.Index()
	oDeletedAtIndexOptions.SetName("game_types-da")

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
	}

	if _, oErr := oIndexView.CreateMany(oAbstractMongodb.Context, aIndexModels); oErr != nil {
		return nil, oErr
	}

	oCounters := oAbstractMongodb.Database.Collection("counters")

	return &GameTypeModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oCounters,
	}, nil
}

func (oSelf *GameTypeModel) nextId() (uint, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "game_type"},
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

func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, bson.M{
		"parent_id":  iParentId,
		"deleted_at": oDeletedAtZero,
	})
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.GameType{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oNew.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oNew.Name = *oValue.Name
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oValue.Key != nil {
		oSet["key"] = *oValue.Key
	}
	if oValue.Name != nil {
		oSet["name"] = *oValue.Name
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

func (oSelf *GameTypeModel) RemoveOneById(iId uint64) error {
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
