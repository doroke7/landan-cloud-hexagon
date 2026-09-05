package outputApplicationMysqlModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	mysqlBase "example/internal/output/application/mysql"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionModel struct {
	*mysqlBase.AbstractMysql
}

func NewAdminPermissionModel(oAbstractModel *mysqlBase.AbstractMysql) outputPortAnyModel.AdminPermissionModel {
	return &AdminPermissionModel{
		AbstractMysql: oAbstractModel,
	}
}

func (oSelf *AdminPermissionModel) AddOne(oAdminPermission *domain.AdminPermissionValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminPermission)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermission{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows inserted")
	}

	return nil
}

func (oSelf *AdminPermissionModel) ShowOneById(iId uint64) (*domain.AdminPermission, error) {
	var oAdminPermission domain.AdminPermission

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAdminPermission, iId).Error; oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oErr
	}

	return &oAdminPermission, nil
}

func (oSelf *AdminPermissionModel) EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminPermission)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermission{}).
		Where("id = ?", iId).
		UpdateColumns(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *AdminPermissionModel) RemoveOneById(iId uint64) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermission{}).
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
}

func (oSelf *AdminPermissionModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

	var aAdminPermissions []*domain.AdminPermission

	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermission{}).
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

	if oErr := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aAdminPermissions).Error; oErr != nil {
		return nil, oErr
	}

	return aAdminPermissions, nil
}

func (oSelf *AdminPermissionModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminPermission{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}
