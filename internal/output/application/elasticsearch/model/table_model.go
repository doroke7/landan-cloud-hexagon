package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	domain "example/internal/domain"
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type TableModel struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
	Index string
}

func NewTableModel(oAbstractModel *outputApplicationElasticsearch.AbstractElasticsearch) outputPortAnyModel.TableModel {
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

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

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

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
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

	return uint64(iTotal), oErr
}

func (oSelf *TableModel) AddOne(oValue *domain.TableValue) error {
	iId, oErr := oSelf.NextId("table")
	if oErr != nil {
		return oErr
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
		return oErr
	}

	return nil
}

func (oSelf *TableModel) EditOneById(oValue *domain.TableValue, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *TableModel) RemoveOneById(iId uint64) error {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows deleted")
	}

	return nil
}
