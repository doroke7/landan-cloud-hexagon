package outputApplicationPostgresqlLogic

import (
	"errors"
	"strings"
	"sync"

	domain "example/internal/domain"
	postgresqlBase "example/internal/output/application/postgresql"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminUserLogic struct {
	*postgresqlBase.AbstractPostgresql
}

func NewAdminUserLogic(oAbstractLogic *postgresqlBase.AbstractPostgresql) outputPortAnyLogic.AdminUserLogic {
	return &AdminUserLogic{
		AbstractPostgresql: oAbstractLogic,
	}
}

func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractPostgresql.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aAdminUsers []*domain.AdminUser
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.
			DB.
			WithContext(oSelf.Context).
			Model(&domain.AdminUser{}).
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

		oFindErr = oQuery.
			Limit(int(*oLimit.Count)).
			Offset(int(*oLimit.Offset)).
			Find(&aAdminUsers).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminUser{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aAdminUsers, 0, oFindErr
	}

	if oCountErr != nil {
		return aAdminUsers, 0, oCountErr
	}

	return aAdminUsers, uint64(iTotal), nil
}

func (oSelf *AdminUserLogic) AddAminUser(oValue *domain.AdminUserValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminUser{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("新增0筆")
	}

	return nil
}
