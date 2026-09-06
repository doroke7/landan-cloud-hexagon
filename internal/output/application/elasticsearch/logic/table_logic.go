package outputApplicationElasticsearchLogic

import (
	"encoding/json"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type TableLogic struct {
	*AbstractLogic
	Index string
}

func NewTableLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractLogic: oAbstractLogic,
		Index:         oAbstractLogic.IndexName("tables"),
	}
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, 0, oErr
	}

	aTables := make([]*domain.Table, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oTable domain.Table
		if oErr := json.Unmarshal(oHit.Source, &oTable); oErr != nil {
			return nil, 0, oErr
		}
		aTables[i] = &oTable
	}

	return aTables, uint64(oResult.Total), nil
}

func (oSelf *TableLogic) ShowTablesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {
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
