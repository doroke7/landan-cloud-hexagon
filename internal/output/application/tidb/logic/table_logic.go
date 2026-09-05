package outputApplicationTidbLogic

import (
	"strings"
	"sync"

	domain "example/internal/domain"
	tidbBase "example/internal/output/application/tidb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type TableLogic struct {
	*tidbBase.AbstractTidb
}

func NewTableLogic(oAbstractLogic *tidbBase.AbstractTidb) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractTidb: oAbstractLogic,
	}
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint, error) {
	aWheres := oSelf.AbstractTidb.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractTidb.SortersToOrders(aSorters)
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

	return aTables, uint(iTotal), nil
}
