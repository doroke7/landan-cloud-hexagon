package resource

import (
	"errors"

	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AppUserLogic struct {
	*resourceBase.AbstractResource
}

func NewAppUserLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyModel.AppUserModel {
	return &AppUserLogic{
		AbstractResource: oAbstractLogic,
	}
}

func (oSelf *AppUserLogic) AddAppUser(oAppUser *domain.AppUser) (*domain.AppUser, error) {

	return &domain.AppUser{
		Id:       1,
		Name:     "11",
		Password: "222222",
	}, nil
}

func (oSelf *AppUserLogic) IncreaseBalance(id uint, amount uint) (bool, error) {

	return true, nil
}

// ShowOneByName 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AppUserLogic) ShowOneByName(sName string) (*domain.AppUser, error) {
	return nil, errors.New("not supported by resource")
}

// ShowOneById 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AppUserLogic) ShowOneById(iId uint) (*domain.AppUser, error) {
	return nil, errors.New("not supported by resource")
}

// AddOne 目前 Resource gRPC service 的 AddAppUser RPC 只帶 name，沒有 password，
// 不夠支撐帳號註冊需要的欄位，先不支援。
func (oSelf *AppUserLogic) AddOne(oAppUser *domain.AppUserValue) (bool, error) {
	return false, errors.New("not supported by resource")
}
