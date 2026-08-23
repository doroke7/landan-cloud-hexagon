package outputApplicationResourceModel

import (
	"errors"

	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
)

type AppUserModel struct {
	*resourceBase.AbstractResource
}

func NewAppUserModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractResource: oAbstractModel,
	}
}

func (oSelf *AppUserModel) AddOne(oAppUser *domain.AppUserValue) (bool, error) {

	oRequest := &pbResourceModel.AppUserAddOneInput{}

	if oAppUser.Name != nil {
		oRequest.Name = *oAppUser.Name
	}

	if _, oErr := oSelf.ResourceModelClient.AppUser.AddAppUser(oSelf.Context, oRequest); oErr != nil {
		return false, oErr
	}

	return true, nil
}

// ShowOneByName 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {
	return nil, errors.New("not supported by resource")
}

// ShowOneById 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {
	return nil, errors.New("not supported by resource")
}

// IncreaseBalance 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AppUserModel) IncreaseBalance(iId uint, iAmount uint) (bool, error) {
	return false, errors.New("not supported by resource")
}
