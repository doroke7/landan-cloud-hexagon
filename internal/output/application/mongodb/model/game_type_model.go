package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	pkg "example/pkg"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
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

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractMongodb.Context, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("game_types-k"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("game_types-n-da"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("game_types-da"),
		},
	}); oErr != nil {
		return nil, oErr
	}

	return &GameTypeModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oAbstractMongodb.Database.Collection("counters"),
	}, nil
}

func (oSelf *GameTypeModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "game_type"},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	var oCounter struct {
		Seq uint `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *GameTypeModel) ShowOneById(iId uint) (*domain.GameType, error) {
	var oDoc domain.GameTypeDocument

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"_id":        iId,
		"deleted_at": oDeletedAtZero,
	}).Decode(&oDoc)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return domain.GameTypeDocumentToGameType(&oDoc), nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aDocs []*domain.GameTypeDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aGameTypes := make([]*domain.GameType, len(aDocs))
	for i, oDoc := range aDocs {
		aGameTypes[i] = domain.GameTypeDocumentToGameType(oDoc)
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.GameTypeDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint) (bool, error) {
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
		return false, oErr
	}

	return oResult.MatchedCount > 0, nil
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint) (bool, error) {
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
	)
	if oErr != nil {
		return false, oErr
	}

	return oResult.ModifiedCount > 0, nil
}
