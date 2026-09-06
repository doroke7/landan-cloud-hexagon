package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
)

// oDeletedAtZero 跟 mysql/resource adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameModel struct {
	*AbstractModel
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、
// name+deleted_at、deleted_at、game_type_id+deleted_at 供查詢用。
func NewGameModel(oAbstractModel *AbstractModel) (outputPortAnyModel.GameModel, error) {
	oCollection := oAbstractModel.Database.Collection("games")

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractModel.Context, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("games-k"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("games-n-da"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("games-da"),
		},
		{
			Keys:    bson.D{{Key: "game_type_id", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("games-gti-da"),
		},
	}); oErr != nil {
		return nil, oErr
	}

	return &GameModel{
		AbstractModel: oAbstractModel,
		Collection:    oCollection,
		Counters:      oAbstractModel.Database.Collection("counters"),
	}, nil
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

func (oSelf *GameModel) AddOne(oGame *domain.GameValue) error {
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

func (oSelf *GameModel) ShowOneById(iId uint64) (*domain.Game, error) {
	var oGame domain.Game

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"_id":        iId,
		"deleted_at": oDeletedAtZero,
	}).Decode(&oGame)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {
	var oGame domain.Game

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"key":        sKey,
		"deleted_at": oDeletedAtZero,
	}).Decode(&oGame)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameValue, iId uint64) error {
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
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
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

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aGames []*domain.Game
	if oErr := oCursor.All(oSelf.Context, &aGames); oErr != nil {
		return nil, oErr
	}

	return aGames, nil
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
