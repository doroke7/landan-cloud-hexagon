package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputApplicationResource "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
)

type AppUserModel struct {
	*outputApplicationResource.AbstractResource
}

func NewAppUserModel(oAbstractModel *outputApplicationResource.AbstractResource) outputPortAnyModel.AppUserModel {
	return &AppUserModel{
		AbstractResource: oAbstractModel,
	}
}

func (oSelf *AppUserModel) AddOne(oAppUser *domain.AppUserValue) error {

	oRequest := &pbResourceModel.AppUserAddOneInput{
		Variable: &pbResourceModel.AppUserVariable{
			Name:     oAppUser.Name,
			Password: oAppUser.Password,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AppUser.AddAppUser(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {

	oResp, oErr := oSelf.ResourceModelClient.AppUser.ShowOneByName(
		oSelf.Context,
		&pbResourceModel.AppUserShowOneByNameInput{Name: sName},
	)
	if oErr != nil {
		return nil, oErr
	}

	oProtoAppUser := oResp.GetAppUser()

	return &domain.AppUser{
		Id:       uint(oProtoAppUser.GetId()),
		Name:     oProtoAppUser.GetName(),
		Password: oProtoAppUser.GetPassword(),
		Balance:  uint(oProtoAppUser.GetBalance()),
	}, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {

	oResp, oErr := oSelf.ResourceModelClient.AppUser.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AppUserShowOneByIdInput{Id: uint64(iId)},
	)
	if oErr != nil {
		return nil, oErr
	}

	oProtoAppUser := oResp.GetAppUser()

	return &domain.AppUser{
		Id:       uint(oProtoAppUser.GetId()),
		Name:     oProtoAppUser.GetName(),
		Password: oProtoAppUser.GetPassword(),
		Balance:  uint(oProtoAppUser.GetBalance()),
	}, nil
}

func (oSelf *AppUserModel) IncreaseBalance(iId uint, iAmount uint64) error {

	_, oErr := oSelf.ResourceModelClient.AppUser.IncreaseBalance(
		oSelf.Context,
		&pbResourceModel.AppUserIncreaseBalanceInput{Id: uint64(iId), Amount: uint64(iAmount)},
	)

	return oErr
}
