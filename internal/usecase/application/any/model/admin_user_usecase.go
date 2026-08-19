package usecase

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
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

func (oSelf *AdminUserUsecase) ShowOneById(iId uint) (*domain.AdminUser, error) {

	// 不需要 太多判斷，直接回傳 資料庫方法
	oAdminUser, err := oSelf.AdminUserModel.ShowOneById(iId)

	return oAdminUser, err
}
