package outputApplicationMysqlModel

import (
	"context"
	"errors"
	"strings"
	"time"

	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	"gorm.io/gorm"

	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserModel struct {
	*AbstractModel
}

func NewAdminUserModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {
	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.MYSQL.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	var oAdminUser domain.AdminUser

	oResult := oSelf.DB.WithContext(oCurrentContext).
		Where("name = ?", sName).
		First(&oAdminUser)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, oResult.Error
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint64) (*domain.AdminUser, error) {

	var oAdminUser domain.AdminUser
	sKey := oSelf.Aop.Key("AdminUserModel.SObI", iId)
	iTtl := oSelf.Aop.Ttl(30 * time.Minute)

	err := oSelf.Aop.Cacheable(sKey, iTtl, &oAdminUser, func() (interface{}, error) {
		oThisContext, cCancel := context.WithTimeout(
			oSelf.Context,
			time.Duration(bootstrap.CONFIG.MYSQL.TIMEOUT)*time.Millisecond,
		)
		defer cCancel()

		var oAdminUser domain.AdminUser

		oResult := oSelf.DB.WithContext(oThisContext).Preload("AdminRoles").First(&oAdminUser, iId)

		if oResult.Error != nil {
			if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
				return nil, errors.New("record not found")
			}
			return nil, oResult.Error
		}
		return oAdminUser, nil
	})

	return &oAdminUser, err
}

func (oSelf *AdminUserModel) RemoveOneById(iId uint64) error {
	sKey := oSelf.Aop.Key("AdminUserModel.SObI", iId)

	oErr := oSelf.Aop.CacheEvict(sKey, func() error {
		oResult := oSelf.DB.WithContext(oSelf.Context).
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

		return nil
	})

	return oErr
}

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserValue, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminUser)
	if oErr != nil {
		return oErr
	}

	sKey := oSelf.Aop.Key("AdminUserModel.SObI", iId)

	oErr = oSelf.Aop.CacheEvict(sKey, func() error {
		oResult := oSelf.DB.WithContext(oSelf.Context).
			Model(&domain.AdminUser{}).
			Where("id = ?", iId).
			UpdateColumns(oColumns)

		if oResult.Error != nil {
			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			return errors.New("0 rows updated")
		}

		return nil
	})

	return oErr
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

	var aAdminUsers []*domain.AdminUser
	var oAdminUser domain.AdminUser

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUser)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}
	oQuery = oQuery.Where("deleted_at = ?", "2038-01-19 03:14:07")

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil || oOrder.Value == nil {
			continue
		}

		sDirection := "ASC"
		if strings.EqualFold(*oOrder.Value, "desc") {
			sDirection = "DESC"
		}

		oQuery = oQuery.Order(*oOrder.Field + " " + sDirection)
	}

	oResult := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aAdminUsers)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64
	var oAdminUser domain.AdminUser

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUser)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}
	oQuery = oQuery.Where("deleted_at = ?", "2038-01-19 03:14:07")

	oResult := oQuery.Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint64(iTotal), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminUser)
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
		return errors.New("0 rows inserted")
	}

	return nil
}
