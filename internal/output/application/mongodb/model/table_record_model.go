package outputApplicationMongodbModel

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

type TableRecordModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：no 唯一、table_id+game_id+no+deleted_at 供查詢用。
func NewTableRecordModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) (outputPortAnyModel.TableRecordModel, error) {
	oCollection := oAbstractMongodb.Database.Collection("table_records")

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractMongodb.Context, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "no", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("table_records-n"),
		},
		{
			Keys: bson.D{
				{Key: "table_id", Value: 1},
				{Key: "game_id", Value: 1},
				{Key: "no", Value: 1},
				{Key: "deleted_at", Value: 1},
			},
			Options: options.Index().SetName("table_records-ti-gi-n-da"),
		},
	}); oErr != nil {
		return nil, oErr
	}

	return &TableRecordModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oAbstractMongodb.Database.Collection("counters"),
	}, nil
}

func (oSelf *TableRecordModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "table_record"},
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

func (oSelf *TableRecordModel) AddOne(oTableRecord *domain.TableRecordValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.TableRecordDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oTableRecord.No != nil {
		oDoc.No = *oTableRecord.No
	}
	if oTableRecord.GameId != nil {
		oDoc.GameId = *oTableRecord.GameId
	}
	if oTableRecord.TableId != nil {
		oDoc.TableId = *oTableRecord.TableId
	}
	if oTableRecord.State != nil {
		oDoc.State = *oTableRecord.State
	}
	if oTableRecord.Text != nil {
		oDoc.Text = string(*oTableRecord.Text)
	}
	if oTableRecord.Image != nil {
		oDoc.Image = string(*oTableRecord.Image)
	}
	if oTableRecord.Result != nil {
		oDoc.Result = string(*oTableRecord.Result)
	}
	if oTableRecord.StartedAt != nil {
		oDoc.StartedAt = *oTableRecord.StartedAt
	}
	if oTableRecord.EndedAt != nil {
		oDoc.EndedAt = *oTableRecord.EndedAt
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableRecordModel) ShowOneById(iId uint) (*domain.TableRecord, error) {
	var oDoc domain.TableRecordDocument

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

	return domain.TableRecordDocumentToTableRecord(&oDoc), nil
}

func (oSelf *TableRecordModel) EditOneById(oTableRecord *domain.TableRecordValue, iId uint) (bool, error) {
	oSet := bson.M{"updated_at": time.Now()}

	if oTableRecord.No != nil {
		oSet["no"] = *oTableRecord.No
	}
	if oTableRecord.GameId != nil {
		oSet["game_id"] = *oTableRecord.GameId
	}
	if oTableRecord.TableId != nil {
		oSet["table_id"] = *oTableRecord.TableId
	}
	if oTableRecord.State != nil {
		oSet["state"] = *oTableRecord.State
	}
	if oTableRecord.Text != nil {
		oSet["text"] = string(*oTableRecord.Text)
	}
	if oTableRecord.Image != nil {
		oSet["image"] = string(*oTableRecord.Image)
	}
	if oTableRecord.Result != nil {
		oSet["result"] = string(*oTableRecord.Result)
	}
	if oTableRecord.StartedAt != nil {
		oSet["started_at"] = *oTableRecord.StartedAt
	}
	if oTableRecord.EndedAt != nil {
		oSet["ended_at"] = *oTableRecord.EndedAt
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

func (oSelf *TableRecordModel) RemoveOneById(iId uint) (bool, error) {
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

func (oSelf *TableRecordModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecord, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aDocs []*domain.TableRecordDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aTableRecords := make([]*domain.TableRecord, len(aDocs))
	for i, oDoc := range aDocs {
		aTableRecords[i] = domain.TableRecordDocumentToTableRecord(oDoc)
	}

	return aTableRecords, nil
}

func (oSelf *TableRecordModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
