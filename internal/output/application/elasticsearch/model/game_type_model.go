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
	aWheres := oSelf.FiltersToWheres(aFilters)

	iTotal, oErr := oSelf.Count(oSelf.Index, aWheres)

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
