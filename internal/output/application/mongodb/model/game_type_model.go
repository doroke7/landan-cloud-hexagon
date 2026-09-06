package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	pkgInput "example/pkg/input"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type GameTypeModel struct {
	*AbstractModel
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、name+deleted_at、deleted_at 供查詢用。
func NewGameTypeModel(oAbstractModel *AbstractModel) (outputPortAnyModel.GameTypeModel, error) {
	oCollection := oAbstractModel.Database.Collection("game_types")

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractModel.Context, []mongo.IndexModel{
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
		AbstractModel: oAbstractModel,
		Collection:    oCollection,
		Counters:      oAbstractModel.Database.Collection("counters"),
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

// ShowOnesByParentId 撈出指定父類型底下、尚未刪除的子類型（給刪除前的擋關檢查用）。
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

func (oSelf *GameTypeModel) ShowOnesByParentId(iParentId uint64) ([]*domain.GameType, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{
		"parent_id":  iParentId,
		"deleted_at": oDeletedAtZero,
	})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aGameTypes []*domain.GameType
	if oErr := oCursor.All(oSelf.Context, &aGameTypes); oErr != nil {
		return nil, oErr
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aGameTypes []*domain.GameType
	if oErr := oCursor.All(oSelf.Context, &aGameTypes); oErr != nil {
		return nil, oErr
	}

	return aGameTypes, nil
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

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) error {
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

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint64) error {
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
