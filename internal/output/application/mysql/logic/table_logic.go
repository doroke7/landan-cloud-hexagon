package mysql

import (
	"strings"
	"sync"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkg "example/pkg"
)

type TableLogic struct {
	*mysqlBase.AbstractMysql
}

func NewTableLogic(oAbstractLogic *mysqlBase.AbstractMysql) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractMysql: oAbstractLogic,
	}
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, int64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aTableRows []*domain.TableRow
	var oTableRow domain.TableRow
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oTableRow).Where("deleted_at = ?", "2038-01-19 03:14:07")
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
			Find(&aTableRows).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oTableRow).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
	}()

	oWaitGroup.Wait()

	aTables := make([]*domain.Table, len(aTableRows))
	for i, oRow := range aTableRows {
		aTables[i] = domain.TableRowToTable(oRow)
	}

	if oFindErr != nil {
		return aTables, 0, oFindErr
	}

	if oCountErr != nil {
		return aTables, 0, oCountErr
	}

	return aTables, iTotal, nil
}
