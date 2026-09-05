package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase struct {
	*AbstractUsecase
	outputPortAnyModel.AdminUserModel
}

func NewAdminUserUsecase(oAminUserRepository outputPortAnyModel.AdminUserModel, oAbstractUsecase *AbstractUsecase) usecasePortAnyModel.AdminUserUsecase {
	return &AdminUserUsecase{
		AbstractUsecase: oAbstractUsecase,
		AdminUserModel:  oAminUserRepository,
	}
}

func (oSelf *AdminUserUsecase) ShowOneByName(sName string) (*domain.AdminUser, error) {

	// 不需要 太多判斷，直接回傳 資料庫方法
	oAdminUser, err := oSelf.AdminUserModel.ShowOneByName(sName)

	return oAdminUser, err
}

func (oSelf *AdminUserUsecase) ShowOneById(iId uint64) (*domain.AdminUser, error) {

	// 不需要 太多判斷，直接回傳 資料庫方法
	oAdminUser, err := oSelf.AdminUserModel.ShowOneById(iId)

	return oAdminUser, err
}

func (oSelf *AdminUserUsecase) AddOne(oAdminUser *domain.AdminUserValue) error {

	return oSelf.AdminUserModel.AddOne(oAdminUser)
}

func (oSelf *AdminUserUsecase) EditOneById(oAdminUser *domain.AdminUserValue, iId uint64) error {

	return oSelf.AdminUserModel.EditOneById(oAdminUser, iId)
}

func (oSelf *AdminUserUsecase) RemoveOneById(iId uint64) error {

	return oSelf.AdminUserModel.RemoveOneById(iId)
}

func (oSelf *AdminUserUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.AdminUserModel.TotalByFilters(aFilters)

	return uint64(iTotal), oErr
}
