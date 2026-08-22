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

type TableRecordModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewTableRecordModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.TableRecordModel {
	return &TableRecordModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("table_records"),
	}
}

func (oSelf *TableRecordModel) ShowOneById(iId uint) (*domain.TableRecord, error) {
	var oTableRecord domain.TableRecord

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oTableRecord)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oTableRecord.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oTableRecord, nil
}

func (oSelf *TableRecordModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecord, error) {
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

	aTableRecords := make([]*domain.TableRecord, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oTableRecord domain.TableRecord
		if oErr := json.Unmarshal(oHit.Source, &oTableRecord); oErr != nil {
			return nil, oErr
		}
		aTableRecords[i] = &oTableRecord
	}

	return aTableRecords, nil
}

func (oSelf *TableRecordModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)

	iTotal, oErr := oSelf.Count(oSelf.Index, aWheres)

	return iTotal, oErr
}

func (oSelf *TableRecordModel) AddOne(oValue *domain.TableRecordValue) (bool, error) {
	iId, oErr := oSelf.NextId("table_record")
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.TableRecord{
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
	if oValue.TableId != nil {
		oDoc.TableId = *oValue.TableId
	}
	if oValue.State != nil {
		oDoc.State = *oValue.State
	}
	if oValue.Text != nil {
		oDoc.Text = *oValue.Text
	}
	if oValue.Image != nil {
		oDoc.Image = *oValue.Image
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

func (oSelf *TableRecordModel) EditOneById(oValue *domain.TableRecordValue, iId uint) (bool, error) {
	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)

	return bOk, oErr
}

func (oSelf *TableRecordModel) RemoveOneById(iId uint) (bool, error) {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)

	return bOk, oErr
}
