package postgresql

import (
	"context"
	"errors"
	"strings"
	"time"

	pkg "example/pkg"

	"gorm.io/gorm"

	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	postgresqlBase "example/internal/output/application/postgresql"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserModel struct {
	*postgresqlBase.AbstractPostgresql
}

func NewAdminUserModel(oAbstractModel *postgresqlBase.AbstractPostgresql) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractPostgresql: oAbstractModel,
	}
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {
	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.POSTGRESQL.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	var oAdminUserRow domain.AdminUserRow

	if err := oSelf.DB.WithContext(oCurrentContext).Where("name = ?", sName).First(&oAdminUserRow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, err
	}

	return domain.AdminUserRowToAdminUser(&oAdminUserRow), nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {

	var oAdminUserRow domain.AdminUserRow
	sKey := oSelf.Aop.Key("AdminUser.SObI", iId)
	iTtl := oSelf.Aop.Ttl(30 * time.Minute)

	err := oSelf.Aop.Cacheable(sKey, iTtl, &oAdminUserRow, func() (interface{}, error) {
		oThisContext, cCancel := context.WithTimeout(
			oSelf.Context,
			time.Duration(bootstrap.CONFIG.POSTGRESQL.TIMEOUT)*time.Millisecond,
		)
		defer cCancel()

		var oAdminUserRow domain.AdminUserRow
		if err := oSelf.DB.WithContext(oThisContext).First(&oAdminUserRow, iId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("資料不存在")
			}
			return nil, err
		}
		return oAdminUserRow, nil
	})

	return domain.AdminUserRowToAdminUser(&oAdminUserRow), err
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractPostgresql.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractPostgresql.PaginationToLimit(oPagination)

	var aAdminUserRows []*domain.AdminUserRow
	var oAdminUserRow domain.AdminUserRow

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUserRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

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

	if oErr := oQuery.
		Limit(int(*oLimit.Count)).
		Offset(int(*oLimit.Offset)).
		Find(&aAdminUserRows).Error; oErr != nil {
		return nil, oErr
	}

	aAdminUsers := make([]*domain.AdminUser, len(aAdminUserRows))
	for i, oAdminUserRow := range aAdminUserRows {
		aAdminUsers[i] = domain.AdminUserRowToAdminUser(oAdminUserRow)
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.AbstractPostgresql.FiltersToWheres(aFilters)

	var iTotal int64
	var oAdminUserRow domain.AdminUserRow

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUserRow)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {
	var oAdminUserRow domain.AdminUserRow

	oColumns, oErr := pkg.StructToMap(oAdminUser)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oAdminUserRow).
		Create(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
