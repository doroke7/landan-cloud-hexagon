package outputApplicationResourceModel

import (
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

func (oSelf *AppUserModel) ShowOneByName(sName string) (*domain.AppUser, error) {

	oResp, oErr := oSelf.ResourceModelClient.AppUser.ShowOneByName(
		oSelf.Context,
		&pbResourceModel.AppUserShowOneByNameInput{Name: sName},
	)
	if oErr != nil {
		return nil, oErr
	}

	return &domain.AppUser{
		Id:       uint(oResp.GetId()),
		Name:     oResp.GetName(),
		Password: oResp.GetPassword(),
		Balance:  uint(oResp.GetBalance()),
	}, nil
}

func (oSelf *AppUserModel) ShowOneById(iId uint) (*domain.AppUser, error) {

	oResp, oErr := oSelf.ResourceModelClient.AppUser.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AppUserShowOneByIdInput{Id: uint32(iId)},
	)
	if oErr != nil {
		return nil, oErr
	}

	return &domain.AppUser{
		Id:       uint(oResp.GetId()),
		Name:     oResp.GetName(),
		Password: oResp.GetPassword(),
		Balance:  uint(oResp.GetBalance()),
	}, nil
}

func (oSelf *AppUserModel) IncreaseBalance(iId uint, iAmount uint) (bool, error) {

	oResp, oErr := oSelf.ResourceModelClient.AppUser.IncreaseBalance(
		oSelf.Context,
		&pbResourceModel.AppUserIncreaseBalanceInput{Id: uint32(iId), Amount: uint32(iAmount)},
	)
	if oErr != nil {
		return false, oErr
	}

	return oResp.GetStatus(), nil
}
