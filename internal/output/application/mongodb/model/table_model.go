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

type TableModel struct {
	*AbstractModel
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、deleted_at、no+deleted_at 供查詢用。
func NewTableModel(oAbstractModel *AbstractModel) (outputPortAnyModel.TableModel, error) {
	oCollection := oAbstractModel.Database.Collection("tables")
	oIndexView := oCollection.Indexes()

	oKeyIndexOptions := options.Index()
	oKeyIndexOptions.SetUnique(true)
	oKeyIndexOptions.SetName("tables-key")

	oDeletedAtIndexOptions := options.Index()
	oDeletedAtIndexOptions.SetName("tables-da")

	oNoDeletedAtIndexOptions := options.Index()
	oNoDeletedAtIndexOptions.SetName("tables-n-da")

	aIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: oKeyIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: oDeletedAtIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "no", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: oNoDeletedAtIndexOptions,
		},
	}

	if _, oErr := oIndexView.CreateMany(oAbstractModel.Context, aIndexModels); oErr != nil {
		return nil, oErr
	}

	oCounters := oAbstractModel.Database.Collection("counters")

	return &TableModel{
		AbstractModel: oAbstractModel,
		Collection:    oCollection,
		Counters:      oCounters,
	}, nil
}

func (oSelf *TableModel) nextId() (uint, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "table"},
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

func (oSelf *TableModel) AddOne(oTable *domain.TableVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.Table{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oTable.No != nil {
		oNew.No = *oTable.No
	}
	if oTable.GameId != nil {
		oNew.GameId = *oTable.GameId
	}
	if oTable.Key != nil {
		oNew.Key = *oTable.Key
	}
	if oTable.State != nil {
		oNew.State = *oTable.State
	}
	if oTable.Description != nil {
		oNew.Description = *oTable.Description
	}
	if oTable.Result != nil {
		oNew.Result = *oTable.Result
	}
	if oTable.StartedAt != nil {
		oNew.StartedAt = *oTable.StartedAt
	}
	if oTable.EndedAt != nil {
		oNew.EndedAt = *oTable.EndedAt
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oTable.No != nil {
		oSet["no"] = *oTable.No
	}
	if oTable.GameId != nil {
		oSet["game_id"] = *oTable.GameId
	}
	if oTable.Key != nil {
		oSet["key"] = *oTable.Key
	}
	if oTable.State != nil {
		oSet["state"] = *oTable.State
	}
	if oTable.Description != nil {
		oSet["description"] = *oTable.Description
	}
	if oTable.Result != nil {
		oSet["result"] = *oTable.Result
	}
	if oTable.StartedAt != nil {
		oSet["started_at"] = *oTable.StartedAt
	}
	if oTable.EndedAt != nil {
		oSet["ended_at"] = *oTable.EndedAt
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

func (oSelf *TableModel) RemoveOneById(iId uint64) error {
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

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
