package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

// oDeletedAtZero 跟 mysql/resource adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

func NewGameModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oAbstractMongodb.Database.Collection("games"),
		Counters:        oAbstractMongodb.Database.Collection("counters"),
	}
}

func (oSelf *GameModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "game"},
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

func (oSelf *GameModel) AddOne(oGame *domain.GameValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.GameDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oGame.GameTypeId != nil {
		oDoc.GameTypeId = *oGame.GameTypeId
	}
	if oGame.Key != nil {
		oDoc.Key = *oGame.Key
	}
	if oGame.Name != nil {
		oDoc.Name = *oGame.Name
	}
	if oGame.Description != nil {
		oDoc.Description = *oGame.Description
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oDoc domain.GameDocument

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

	return domain.GameDocumentToGame(&oDoc), nil
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameValue, iId uint) (bool, error) {
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
		return false, oErr
	}

	return oResult.MatchedCount > 0, nil
}

func (oSelf *GameModel) RemoveOneById(iId uint) (bool, error) {
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

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersToFindOptions(aSorters)

	iSize := uint(10)
	iPage := uint(1)
	if oPagination != nil {
		if oPagination.Size != nil && *oPagination.Size != 0 {
			iSize = *oPagination.Size
		}
		if oPagination.Page != nil && *oPagination.Page != 0 {
			iPage = *oPagination.Page
		}
	}
	oFindOptions.SetLimit(int64(iSize)).SetSkip(int64((iPage - 1) * iSize))

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aDocs []*domain.GameDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aGames := make([]*domain.Game, len(aDocs))
	for i, oDoc := range aDocs {
		aGames[i] = domain.GameDocumentToGame(oDoc)
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
