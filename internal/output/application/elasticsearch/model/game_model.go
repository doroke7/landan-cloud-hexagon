package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

// oDeletedAtZero 跟 mysql/mongodb adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewGameModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("games"),
	}
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oGame domain.Game

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oGame)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oGame.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {
	oBody := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{"term": map[string]any{"key": sKey}},
					{"term": map[string]any{"deleted_at": oDeletedAtZero}},
				},
			},
		},
	}

	aBodyBytes, oErr := json.Marshal(oBody)
	if oErr != nil {
		return nil, oErr
	}

	aOptions := []func(*esapi.SearchRequest){
		oSelf.Client.Search.WithContext(oSelf.Context),
		oSelf.Client.Search.WithIndex(oSelf.Index),
		oSelf.Client.Search.WithBody(bytes.NewReader(aBodyBytes)),
		oSelf.Client.Search.WithSize(1),
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	if len(oResult.Hits) == 0 {
		return nil, nil
	}

	var oGame domain.Game
	if oErr := json.Unmarshal(oResult.Hits[0].Source, &oGame); oErr != nil {
		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOnesByGameTypeId(iGameTypeId uint) ([]*domain.Game, error) {
	sField := "game_type_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iGameTypeId},
	}

	iSize := uint(10000)
	iPage := uint(1)
	aGames, oErr := oSelf.ShowOnesByFiltersWithOrdersPagination(aFilters, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
	return aGames, oErr
}

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint) (uint, error) {
	sGameTypeIdField := "game_type_id"
	sDeletedAtField := "deleted_at"
	aFilters := []*pkgInput.Filter{
		{Field: &sGameTypeIdField, Value: iGameTypeId},
		{Field: &sDeletedAtField, Value: oDeletedAtZero},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
}

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {
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

	aGames := make([]*domain.Game, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oGame domain.Game
		if oErr := json.Unmarshal(oHit.Source, &oGame); oErr != nil {
			return nil, oErr
		}
		aGames[i] = &oGame
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint, error) {
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

	return uint(iTotal), oErr
}

func (oSelf *GameModel) AddOne(oValue *domain.GameValue) error {
	iId, oErr := oSelf.NextId("game")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.Game{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.GameTypeId != nil {
		oDoc.GameTypeId = *oValue.GameTypeId
	}
	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}
	if oValue.Description != nil {
		oDoc.Description = *oValue.Description
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *GameModel) EditOneById(oValue *domain.GameValue, iId uint) error {
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
		return errors.New("更新0筆")
	}

	return nil
}

func (oSelf *GameModel) RemoveOneById(iId uint) error {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("刪除0筆")
	}

	return nil
}
