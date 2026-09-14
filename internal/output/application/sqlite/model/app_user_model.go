package outputApplicationSqliteModel

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
	oExpr := gorm.Expr("balance + ?", amount)

	err := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).
		Where("id = ?", id).
		UpdateColumn("balance", oExpr).Error

	if err != nil {
		return err
	}

	return nil
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	oErr := oSelf.DB.WithContext(oSelf.Context).Where("name = ?", sName).First(&oAppUser).Error

	if oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			oNotFoundError := errors.New("record not found")

			return nil, oNotFoundError
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	var oAppUser domain.AppUser

	oErr := oSelf.DB.WithContext(oSelf.Context).First(&oAppUser, iId).Error

	if oErr != nil {
		if errors.Is(oErr, gorm.ErrRecordNotFound) {
			oNotFoundError := errors.New("record not found")

			return nil, oNotFoundError
		}
		return nil, oErr
	}

	return &oAppUser, nil
}

func (oSelf *AppUserModel) AddOne(oValue *domain.AppUserValue) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	oResult := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AppUser{}).Create(oColumns)
	if oResult.Error != nil {
		return oResult.Error
	}

	if oResult.RowsAffected == 0 {
		oZeroRowsError := errors.New("0 rows inserted")

		return oZeroRowsError
	}

	return nil
}
