package outputApplicationMysqlLogic

import (
	"errors"
	"strings"
	"sync"

	"gorm.io/gorm"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminUserLogic struct {
	*mysqlBase.AbstractMysql
}

func NewAdminUserLogic(oAbstractLogic *mysqlBase.AbstractMysql) outputPortAnyLogic.AdminUserLogic {
	return &AdminUserLogic{
		AbstractMysql: oAbstractLogic,
	}
}

func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
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

	delete(oColumns, "admin_role_ids")

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		oResult := oTx.
			Model(&domain.AdminUser{}).
			Create(oColumns)

		if oResult.Error != nil {
			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			return errors.New("新增0筆")
		}

		var iAdminUserId uint64
		if oErr := oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminUserId).Error; oErr != nil {
			return oErr
		}

		if oErr := oTx.Where("admin_user_id = ?", iAdminUserId).Delete(&domain.AdminUsersToAdminRole{}).Error; oErr != nil {
			return oErr
		}

		if len(oValue.AdminRoleIds) == 0 {
			return nil // 沒有角色也算成功，直接結束，不會再往下插入
		}

		aRelations := make([]domain.AdminUsersToAdminRole, 0, len(oValue.AdminRoleIds))
		for _, iAdminRoleId := range oValue.AdminRoleIds {
			aRelations = append(aRelations, domain.AdminUsersToAdminRole{
				AdminUserId: uint(iAdminUserId),
				AdminRoleId: uint(iAdminRoleId),
			})
		}

		oErr = oTx.Create(&aRelations).Error

		return oErr
	})

	return oError
}
