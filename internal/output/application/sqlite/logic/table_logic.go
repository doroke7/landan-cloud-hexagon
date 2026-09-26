package outputApplicationSqliteLogic

import (
	"errors"
	"strings"
	"sync"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationSqlite "example/internal/output/application/sqlite"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type TableLogic struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewTableLogic(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractSqlite: oAbstractSqlite,
	}
}

func (oSelf *TableLogic) ShowTableById(iId uint64) (*domain.Table, error) {
	var oTable domain.Table

	oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&oTable).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oTable, iId).Error

	if oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oTable, nil
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aTables []*domain.Table
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.Table{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		for _, oOrder := range aOrders {
			if oOrder == nil || oOrder.Field == nil {
				continue
			}

			sDirection := "ASC"
			if oOrder.Value != nil && strings.EqualFold(*oOrder.Value, "desc") {
				sDirection = "DESC"
			}

			oQuery = oQuery.Order(*oOrder.Field + " " + sDirection)
		}

		oFindErr = oQuery.
			Limit(int(*oLimit.Count)).
			Offset(int(*oLimit.Offset)).
			Find(&aTables).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.Table{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aTables, 0, oFindErr
	}

	if oCountErr != nil {
		return aTables, 0, oCountErr
	}

	return aTables, uint64(iTotal), nil
}

func (oSelf *TableLogic) ShowTablesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractSqlite.PaginationToLimit(oPagination)

	var aTables []*domain.Table

	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Preload("Game").
		Preload("Game.GameType").
		Model(&domain.Table{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07")

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil {
			continue
		}

		sDirection := "ASC"
		if oOrder.Value != nil && strings.EqualFold(*oOrder.Value, "desc") {
			sDirection = "DESC"
		}

		oQuery = oQuery.Order(*oOrder.Field + " " + sDirection)
	}

	oErr := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aTables).Error

	if oErr != nil {
		return nil, oErr
	}

	return aTables, nil
}
