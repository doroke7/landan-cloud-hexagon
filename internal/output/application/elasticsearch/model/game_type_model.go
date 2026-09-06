package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	domain "example/internal/domain"
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type GameTypeModel struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
	Index string
}

func NewGameTypeModel(oAbstractModel *outputApplicationElasticsearch.AbstractElasticsearch) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("game_types"),
	}
}

func (oSelf *GameTypeModel) ShowOneById(iId uint64) (*domain.GameType, error) {
	var oGameType domain.GameType

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oGameType)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oGameType.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oGameType, nil
}

func (oSelf *GameTypeModel) ShowOnes() ([]*domain.GameType, error) {
	iSize := uint(10000)
	iPage := uint(1)
	return oSelf.ShowOnesByFiltersWithSortersPagination(nil, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
}

// ShowOnesByParentId 撈出指定父類型底下、尚未刪除的子類型（給刪除前的擋關檢查用）。
// 走既有的 filters 查詢，deleted_at 的過濾由 ShowOnesByFiltersWithSortersPagination 內部補上。
func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	sParentIdField := "parent_id"
	sDeletedAtField := "deleted_at"
	aFilters := []*pkgInput.Filter{
		{Field: &sParentIdField, Value: iParentId},
		{Field: &sDeletedAtField, Value: oDeletedAtZero},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
}

func (oSelf *GameTypeModel) ShowOnesByParentId(iParentId uint64) ([]*domain.GameType, error) {
	sField := "parent_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iParentId},
	}

	iSize := uint(10000)
	iPage := uint(1)
	return oSelf.ShowOnesByFiltersWithSortersPagination(aFilters, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {
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

	aGameTypes := make([]*domain.GameType, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oGameType domain.GameType
		if oErr := json.Unmarshal(oHit.Source, &oGameType); oErr != nil {
			return nil, oErr
		}
		aGameTypes[i] = &oGameType
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
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

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) error {
	iId, oErr := oSelf.NextId("game_type")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.GameType{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint64) error {
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

func (oSelf *GameTypeModel) RemoveOneById(iId uint64) error {
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
