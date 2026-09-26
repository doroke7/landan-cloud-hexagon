package outputApplicationMysqlLogic

import (
	"errors"
	"strings"
	"sync"
	"time"

	domain "example/internal/domain"
	outputApplicationMysql "example/internal/output/application/mysql"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	"gorm.io/gorm"
)

type AdminUserLogic struct {
	*outputApplicationMysql.AbstractMysql
}

func NewAdminUserLogic(oAbstractMysql *outputApplicationMysql.AbstractMysql) outputPortAnyLogic.AdminUserLogic {
	oLogic := &AdminUserLogic{
		AbstractMysql: oAbstractMysql,
	}

	return oLogic
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
			Preload("AdminRoles").
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

		oResult := oQuery.
			Limit(int(*oLimit.Count)).
			Offset(int(*oLimit.Offset)).
			Find(&aAdminUsers)
		oFindErr = oResult.Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).
			Model(&domain.AdminUser{}).
			Where("deleted_at = ?", "2038-01-19 03:14:07")

		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oResult := oQuery.Count(&iTotal)
		oCountErr = oResult.Error
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

	if len(oColumns) == 0 {
		return errors.New("0 columns to insert")

	}

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		oResult := oTx.
			Model(&domain.AdminUser{}).
			Create(oColumns)

		if oResult.Error != nil {
			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			return errors.New("0 rows inserted")
		}

		var iAdminUserId uint64
		oResult = oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminUserId)

		if oResult.Error != nil {
			return oResult.Error
		}

		oResult = oTx.Where("admin_user_id = ?", iAdminUserId).Delete(&domain.AdminUsersToAdminRole{})

		if oResult.Error != nil {
			return oResult.Error
		}

		if len(oValue.AdminRoleIds) == 0 {
			return nil // 沒有角色也算成功，直接結束，不會再往下插入
		}

		aRelations := make([]domain.AdminUsersToAdminRole, 0, len(oValue.AdminRoleIds))
		for _, iAdminRoleId := range oValue.AdminRoleIds {
			oRelation := domain.AdminUsersToAdminRole{
				AdminUserId: iAdminUserId,
				AdminRoleId: iAdminRoleId,
			}
			aRelations = append(aRelations, oRelation)
		}

		oResult = oTx.Create(&aRelations)
		oErr = oResult.Error

		return oErr
	})

	return oError
}

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_role_ids")

	sKey := oSelf.Aop.Key("AdminUserLogic.AUBI", iId)

	oErr = oSelf.Aop.CacheEvict(sKey, func() error {
		oTxErr := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

			if len(oColumns) > 0 {
				oResult := oTx.
					Model(&domain.AdminUser{}).
					Where("id = ?", iId).
					UpdateColumns(oColumns)

				if oResult.Error != nil {
					return oResult.Error
				}

				if oResult.RowsAffected == 0 {
					return errors.New("0 rows updated")
				}
			}

			oResult := oTx.Where("admin_user_id = ?", iId).Delete(&domain.AdminUsersToAdminRole{})

			if oResult.Error != nil {
				return oResult.Error
			}

			if len(oValue.AdminRoleIds) == 0 {
				return nil // 沒有角色也算成功，直接結束，不會再往下插入
			}

			aRelations := make([]domain.AdminUsersToAdminRole, 0, len(oValue.AdminRoleIds))
			for _, iAdminRoleId := range oValue.AdminRoleIds {
				oRelation := domain.AdminUsersToAdminRole{
					AdminUserId: iId,
					AdminRoleId: iAdminRoleId,
				}
				aRelations = append(aRelations, oRelation)
			}

			oResult = oTx.Create(&aRelations)
			oError := oResult.Error

			return oError
		})

		return oTxErr
	})

	return oErr
}

func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	sKey := oSelf.Aop.Key("AdminUserLogic.AUBI", iId)
	iTtl := oSelf.Aop.Ttl(30 * time.Minute)

	oErr := oSelf.Aop.Cacheable(sKey, iTtl, &oAdminUser, func() (interface{}, error) {
		var oFresh domain.AdminUser

		oResult := oSelf.DB.WithContext(oSelf.Context).
			Preload("AdminRoles").
			Where("deleted_at = ?", "2038-01-19 03:14:07").
			First(&oFresh, iId)

		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}

		if oResult.Error != nil {
			return nil, oResult.Error
		}

		return oFresh, nil
	})

	if oErr != nil {
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {

	sKey := oSelf.Aop.Key("AdminUserLogic.AUBI", iId)

	oErr := oSelf.Aop.CacheEvict(sKey, func() error {
		oTxErr := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

			oResult := oTx.
				Model(&domain.AdminUser{}).
				Where("id = ?", iId).
				Where("deleted_at = ?", "2038-01-19 03:14:07").
				UpdateColumn("deleted_at", time.Now())

			if oResult.Error != nil {
				return oResult.Error
			}

			if oResult.RowsAffected == 0 {
				return errors.New("0 rows deleted")
			}

			oResult = oTx.
				Where("admin_user_id = ?", iId).
				Delete(&domain.AdminUsersToAdminRole{})
			oError := oResult.Error

			return oError
		})

		return oTxErr
	})

	return oErr
}
