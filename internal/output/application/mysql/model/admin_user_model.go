package mysql

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
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
		time.Duration(bootstrap.CONFIG.DATABASE.TIMEOUT)*time.Millisecond,
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
			time.Duration(bootstrap.CONFIG.DATABASE.TIMEOUT)*time.Millisecond,
		)
		defer cCancel()

		var oDbAdminUser domain.AdminUser
		if err := oSelf.DB.WithContext(oThisContext).First(&oDbAdminUser, iId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("資料不存在")
			}
			return nil, err
		}
		return oDbAdminUser, nil
	})

	return &oAdminUser, err
}

func (oSelf *AdminUserModel) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.AdminUser, error) {
	var aAdminUsers []*domain.AdminUser

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminUser{})

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

func (oSelf *AdminUserModel) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {
	var iTotal int64

	oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminUser{})

	for _, oWhere := range aWheres {
		oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
	}

	if oErr := oQuery.Count(&iTotal).Error; oErr != nil {
		return 0, oErr
	}

	return uint64(iTotal), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {

	oColumns, oErr := pkg.StructToMap(oAdminUser)
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
