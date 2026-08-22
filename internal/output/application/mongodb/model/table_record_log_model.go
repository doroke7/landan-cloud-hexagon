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

type TableRecordLogModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：table_record_id+game_id+state 供查詢用。
func NewTableRecordLogModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) (outputPortAnyModel.TableRecordLogModel, error) {
	oCollection := oAbstractMongodb.Database.Collection("table_record_logs")

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractMongodb.Context, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "table_record_id", Value: 1},
				{Key: "game_id", Value: 1},
				{Key: "state", Value: 1},
			},
			Options: options.Index().SetName("table_record_logs-tri-gi-s"),
		},
	}); oErr != nil {
		return nil, oErr
	}

	return &TableRecordLogModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oAbstractMongodb.Database.Collection("counters"),
	}, nil
}

func (oSelf *TableRecordLogModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "table_record_log"},
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

func (oSelf *TableRecordLogModel) AddOne(oTableRecordLog *domain.TableRecordLogValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.TableRecordLogDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oTableRecordLog.GameId != nil {
		oDoc.GameId = *oTableRecordLog.GameId
	}
	if oTableRecordLog.TableRecordId != nil {
		oDoc.TableRecordId = *oTableRecordLog.TableRecordId
	}
	if oTableRecordLog.State != nil {
		oDoc.State = *oTableRecordLog.State
	}
	if oTableRecordLog.Text != nil {
		oDoc.Text = string(*oTableRecordLog.Text)
	}
	if oTableRecordLog.Image != nil {
		oDoc.Image = string(*oTableRecordLog.Image)
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableRecordLogModel) ShowOneById(iId uint) (*domain.TableRecordLog, error) {
	var oDoc domain.TableRecordLogDocument

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

	return domain.TableRecordLogDocumentToTableRecordLog(&oDoc), nil
}

func (oSelf *TableRecordLogModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecordLog, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aDocs []*domain.TableRecordLogDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aTableRecordLogs := make([]*domain.TableRecordLog, len(aDocs))
	for i, oDoc := range aDocs {
		aTableRecordLogs[i] = domain.TableRecordLogDocumentToTableRecordLog(oDoc)
	}

	return aTableRecordLogs, nil
}

func (oSelf *TableRecordLogModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
