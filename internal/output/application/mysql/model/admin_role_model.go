package outputApplicationMysqlModel

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminRoleModel struct {
	*AbstractModel
}

func NewAdminRoleModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminRoleModel {
	return &AdminRoleModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AdminRoleModel) AddOne(oAdminRole *domain.AdminRoleVariable) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminRole)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows inserted")
	}

	return nil
}

func (oSelf *AdminRoleModel) ShowOneById(iId uint64) (*domain.AdminRole, error) {
	var oAdminRole domain.AdminRole

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAdminRole, iId)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oResult.Error
	}

	return &oAdminRole, nil
}

func (oSelf *AdminRoleModel) ShowOnes() ([]*domain.AdminRole, error) {
	var aAdminRoles []*domain.AdminRole

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aAdminRoles)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aAdminRoles, nil
}

func (oSelf *AdminRoleModel) EditOneById(oAdminRole *domain.AdminRoleVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oAdminRole)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
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

func (oSelf *AdminRoleModel) RemoveOneById(iId uint64) error {
	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
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

func (oSelf *AdminRoleModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractMysql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractMysql.PaginationToLimit(oPagination)

	var aAdminRoles []*domain.AdminRole

	oQuery := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminRole{}).
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
		Find(&aAdminRoles)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aAdminRoles, nil
}

func (oSelf *AdminRoleModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractMysql.FiltersToWheres(aFilters)

	var iTotal int64

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminRole{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	oResult := oQuery.Count(&iTotal)

	if oResult.Error != nil {
		return 0, oResult.Error
	}

	return uint64(iTotal), nil
}
