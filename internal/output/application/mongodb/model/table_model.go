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

type TableModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

func NewTableModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oAbstractMongodb.Database.Collection("tables"),
		Counters:        oAbstractMongodb.Database.Collection("counters"),
	}
}

func (oSelf *TableModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "table"},
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

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.TableDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oTable.No != nil {
		oDoc.No = *oTable.No
	}
	if oTable.GameId != nil {
		oDoc.GameId = *oTable.GameId
	}
	if oTable.Key != nil {
		oDoc.Key = *oTable.Key
	}
	if oTable.State != nil {
		oDoc.State = *oTable.State
	}
	if oTable.Description != nil {
		oDoc.Description = *oTable.Description
	}
	if oTable.Result != nil {
		oDoc.Result = string(*oTable.Result)
	}
	if oTable.StartedAt != nil {
		oDoc.StartedAt = *oTable.StartedAt
	}
	if oTable.EndedAt != nil {
		oDoc.EndedAt = *oTable.EndedAt
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {
	var oDoc domain.TableDocument

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

	return domain.TableDocumentToTable(&oDoc), nil
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint) (bool, error) {
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
		oSet["result"] = string(*oTable.Result)
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
		return false, oErr
	}

	return oResult.MatchedCount > 0, nil
}

func (oSelf *TableModel) RemoveOneById(iId uint) (bool, error) {
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

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error) {
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

	var aDocs []*domain.TableDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aTables := make([]*domain.Table, len(aDocs))
	for i, oDoc := range aDocs {
		aTables[i] = domain.TableDocumentToTable(oDoc)
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}
