package outputApplicationMysqlModel

import (
	"errors"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgUtility "example/pkg/utility"
)

type AppUserModel struct {
	*AbstractModel
}

func NewAppUserModel(oAbstractModel *AbstractModel) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *AppUserModel) IncreaseBalance(id uint, amount uint64) error {

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AppUser{}).
		Where("id = ?", id).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount))

	if oResult.Error != nil {
		return oResult.Error
	}

	return nil
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Where("name = ?", sName).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAppUser)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, oResult.Error
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAppUser, iId)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, oResult.Error
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) error {

	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AppUser{}).
		Create(oColumns)

	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		return errors.New("0 rows inserted")
	}

	return nil
}
