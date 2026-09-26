package outputApplicationSqliteLogic

import (
	"errors"
	"strings"
	"sync"
	"time"

	domain "example/internal/domain"
	outputApplicationSqlite "example/internal/output/application/sqlite"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	"gorm.io/gorm"
)

type AdminUserLogic struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewAdminUserLogic(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) outputPortAnyLogic.AdminUserLogic {
	return &AdminUserLogic{
		AbstractSqlite: oAbstractSqlite,
	}
}

func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
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

func (oSelf *AdminUserLogic) AddAdminUser(oValue *domain.AdminUserVariable) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_role_ids")

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminUser{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows inserted")

		return oZeroRowsError
	}

	return nil
}

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_role_ids")

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminUser{}).
		Where("id = ?", iId).
		UpdateColumns(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows updated")

		return oZeroRowsError
	}

	return nil
}

func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("AdminRoles").
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAdminUser, iId).Error

	if errors.Is(oErr, gorm.ErrRecordNotFound) {
		oNotFoundError := errors.New("record not found")

		return nil, oNotFoundError
	}

	if oErr != nil {
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {
	oNow := time.Now()

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminUser{}).
		Where("id = ?", iId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		UpdateColumn("deleted_at", oNow)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows deleted")

		return oZeroRowsError
	}

	return nil
}
