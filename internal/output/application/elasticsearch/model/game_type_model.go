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

type GameTypeModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewGameTypeModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("game_types"),
	}
}

func (oSelf *GameTypeModel) ShowOneById(iId uint) (*domain.GameType, error) {
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

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error) {
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

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := make([]map[string]any, 0, len(aFilters))

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
			aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"term": map[string]any{sField: oValue}}}})
		case "gt", "gte", "lt", "lte":
			aWheres = append(aWheres, map[string]any{"range": map[string]any{sField: map[string]any{sOperator: oValue}}})
		case "contains":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}})
			}
		case "notContains":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}}}})
			}
		case "startsWith":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"prefix": map[string]any{sField: map[string]any{"value": sValue, "case_insensitive": true}}})
			}
		case "endsWith":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue, "case_insensitive": true}}})
			}
		case "in":
			aWheres = append(aWheres, map[string]any{"terms": map[string]any{sField: oValue}})
		case "notIn":
			aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"terms": map[string]any{sField: oValue}}}})
		case "between":
			if aRange, bOk := oValue.([]any); bOk && len(aRange) == 2 {
				aWheres = append(aWheres, map[string]any{"range": map[string]any{sField: map[string]any{"gte": aRange[0], "lte": aRange[1]}}})
			}
		case "match":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"match": map[string]any{sField: sValue}})
			}
		default:
			aWheres = append(aWheres, map[string]any{"term": map[string]any{sField: oValue}})
		}
	}

	oQuery := map[string]any{"match_all": map[string]any{}}
	if len(aWheres) > 0 {
		oQuery = map[string]any{"bool": map[string]any{"filter": aWheres}}
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

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) (bool, error) {
	iId, oErr := oSelf.NextId("game_type")
	if oErr != nil {
		return false, oErr
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
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint) (bool, error) {
	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)

	return bOk, oErr
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint) (bool, error) {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)

	return bOk, oErr
}
