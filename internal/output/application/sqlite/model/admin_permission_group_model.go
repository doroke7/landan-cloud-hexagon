package outputApplicationSqliteModel

import (
	"errors"
	"strings"
	"time"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionGroupModel struct {
	*AbstractModel
}

func NewAdminPermissionGroupModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminPermissionGroupModel {
	return &AdminPermissionGroupModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByParentId(iParentId uint64) ([]*domain.AdminPermissionGroup, error) {
	var aAdminPermissionGroups []*domain.AdminPermissionGroup

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Preload("Parent").
		Preload("Children").
		Model(&domain.AdminPermissionGroup{}).
		Where("parent_id = ?", iParentId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aAdminPermissionGroups)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractSqlite.PaginationToLimit(oPagination)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup

	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Preload("Parent").
		Preload("Children").
		Model(&domain.AdminPermissionGroup{}).
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
		Find(&aAdminPermissionGroups).Error

	if oErr != nil {
		return nil, oErr
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)

	var iTotal int64
	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	oErr := oQuery.Count(&iTotal).Error

	if oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *AdminPermissionGroupModel) AddOne(oValue *domain.AdminPermissionGroupVariable) error {
	oColumns, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
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

func (oSelf *AdminPermissionGroupModel) EditOneById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oColumns, _ := pkgUtility.StructToMap(oValue)

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
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

func (oSelf *AdminPermissionGroupModel) RemoveOneById(iId uint64) error {
	oNow := time.Now()

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
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
