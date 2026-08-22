package elasticsearch

import (
	"encoding/json"
	"strconv"
	"time"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type TableModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewTableModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("tables"),
	}
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {
	var oTable domain.Table

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oTable)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oTable.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oTable, nil
}

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkg.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aTables := make([]*domain.Table, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oTable domain.Table
		if oErr := json.Unmarshal(oHit.Source, &oTable); oErr != nil {
			return nil, oErr
		}
		aTables[i] = &oTable
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)

	iTotal, oErr := oSelf.Count(oSelf.Index, aWheres)

	return iTotal, oErr
}

func (oSelf *TableModel) AddOne(oValue *domain.TableValue) (bool, error) {
	iId, oErr := oSelf.NextId("table")
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.Table{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.No != nil {
		oDoc.No = *oValue.No
	}
	if oValue.GameId != nil {
		oDoc.GameId = *oValue.GameId
	}
	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.State != nil {
		oDoc.State = *oValue.State
	}
	if oValue.Description != nil {
		oDoc.Description = *oValue.Description
	}
	if oValue.Result != nil {
		oDoc.Result = *oValue.Result
	}
	if oValue.StartedAt != nil {
		oDoc.StartedAt = *oValue.StartedAt
	}
	if oValue.EndedAt != nil {
		oDoc.EndedAt = *oValue.EndedAt
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableModel) EditOneById(oValue *domain.TableValue, iId uint) (bool, error) {
	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)

	return bOk, oErr
}

func (oSelf *TableModel) RemoveOneById(iId uint) (bool, error) {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)

	return bOk, oErr
}
