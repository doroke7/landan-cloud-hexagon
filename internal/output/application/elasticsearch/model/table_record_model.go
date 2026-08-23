package elasticsearch

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

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
	aFilterClauses := make([]map[string]any, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		sField := *oFilter.Field
		oValue := oFilter.Value

		sOperator := "eq"
		if oFilter.Operator != nil {
			sOperator = *oFilter.Operator
		}

		switch sOperator {
		case "ne":
			aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"term": map[string]any{sField: oValue}}}})
		case "gt", "gte", "lt", "lte":
			aFilterClauses = append(aFilterClauses, map[string]any{"range": map[string]any{sField: map[string]any{sOperator: oValue}}})
		case "contains":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}})
			}
		case "notContains":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}}}})
			}
		case "startsWith":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"prefix": map[string]any{sField: map[string]any{"value": sValue, "case_insensitive": true}}})
			}
		case "endsWith":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue, "case_insensitive": true}}})
			}
		case "in":
			aFilterClauses = append(aFilterClauses, map[string]any{"terms": map[string]any{sField: oValue}})
		case "notIn":
			aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"terms": map[string]any{sField: oValue}}}})
		case "between":
			if aRange, bOk := oValue.([]any); bOk && len(aRange) == 2 {
				aFilterClauses = append(aFilterClauses, map[string]any{"range": map[string]any{sField: map[string]any{"gte": aRange[0], "lte": aRange[1]}}})
			}
		case "match":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"match": map[string]any{sField: sValue}})
			}
		default:
			aFilterClauses = append(aFilterClauses, map[string]any{"term": map[string]any{sField: oValue}})
		}
	}

	oQuery := map[string]any{"match_all": map[string]any{}}
	if len(aFilterClauses) > 0 {
		oQuery = map[string]any{"bool": map[string]any{"filter": aFilterClauses}}
	}

	aBodyBytes, oErr := json.Marshal(map[string]any{"query": oQuery})
	if oErr != nil {
		return 0, oErr
	}

	aOptions := []func(*esapi.CountRequest){
		oSelf.Client.Count.WithContext(oSelf.Context),
		oSelf.Client.Count.WithIndex(oSelf.Index),
		oSelf.Client.Count.WithBody(bytes.NewReader(aBodyBytes)),
	}

	iTotal, oErr := oSelf.CountWithOptions(aOptions)

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
