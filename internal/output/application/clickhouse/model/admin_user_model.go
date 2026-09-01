package outputApplicationClickhouseModel

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
	clickhouseBase "example/internal/output/application/clickhouse"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserModel struct {
	*clickhouseBase.AbstractClickhouse
}

func NewAdminUserModel(oAbstractModel *clickhouseBase.AbstractClickhouse) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractClickhouse: oAbstractModel,
	}
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {
	oCurrentContext, cCancel := context.WithTimeout(
		oSelf.Context,
		time.Duration(bootstrap.CONFIG.CLICKHOUSE.TIMEOUT)*time.Millisecond,
	)
	defer cCancel()

	var oAdminUser domain.AdminUser

	if err := oSelf.DB.WithContext(oCurrentContext).Where("name = ?", sName).First(&oAdminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("資料不存在")
		}
		return nil, err
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {

	var oAdminUser domain.AdminUser
	sKey := oSelf.Aop.Key("AdminUser.SObI", iId)
	iTtl := oSelf.Aop.Ttl(30 * time.Minute)

	err := oSelf.Aop.Cacheable(sKey, iTtl, &oAdminUser, func() (interface{}, error) {
		oThisContext, cCancel := context.WithTimeout(
			oSelf.Context,
			time.Duration(bootstrap.CONFIG.CLICKHOUSE.TIMEOUT)*time.Millisecond,
		)
		defer cCancel()

		var oAdminUser domain.AdminUser
		if err := oSelf.DB.WithContext(oThisContext).First(&oAdminUser, iId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("資料不存在")
			}
			return nil, err
		}
		return oAdminUser, nil
	})

	return &oAdminUser, err
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error) {
	aWheres := oSelf.AbstractClickhouse.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractClickhouse.SortersToOrders(aSorters)
	oLimit := oSelf.AbstractClickhouse.PaginationToLimit(oPagination)

	var aAdminUsers []*domain.AdminUser
	var oAdminUser domain.AdminUser

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUser)

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
		Find(&aAdminUsers).Error; oErr != nil {
		return nil, oErr
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aWheres := oSelf.AbstractClickhouse.FiltersToWheres(aFilters)

	var iTotal int64
	var oAdminUser domain.AdminUser

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&oAdminUser)

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {
	oColumns, oErr := pkgUtility.StructToMap(oAdminUser)
	if oErr != nil {
		return false, oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminUser{}).
		Create(oColumns)

	if oResult.Error != nil {
		return false, oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return false, errors.New("新增0筆")
	}

	return true, nil
}
