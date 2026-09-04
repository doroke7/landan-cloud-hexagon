package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

func protoAdminUserToDomainAdminUser(oProto *pbResource.AdminUser) domain.AdminUser {
	if oProto == nil {
		return domain.AdminUser{}
	}

	return domain.AdminUser{
		Id:        uint(oProto.GetId()),
		Name:      oProto.GetName(),
		Password:  oProto.GetPassword(),
		CreatedAt: oProto.GetCreatedAt().AsTime(),
		UpdatedAt: oProto.GetUpdatedAt().AsTime(),
		DeletedAt: oProto.GetDeletedAt().AsTime(),
	}
}

type AdminUserModel struct {
	*resourceBase.AbstractResource
}

func NewAdminUserModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractResource: oAbstractModel,
	}
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {

	oResp, err := oSelf.ResourceModelClient.AdminUser.ShowOneByName(
		oSelf.Context,
		&pbResourceModel.AdminUserShowOneByNameInput{Name: sName},
	)

	if err != nil {
		return nil, err
	}

	oAdminUser := protoAdminUserToDomainAdminUser(oResp.GetAdminUser())

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {

	oResp, err := oSelf.ResourceModelClient.AdminUser.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminUserShowOneByIdInput{Id: uint32(iId)},
	)

	if err != nil {
		return nil, err
	}

	oAdminUser := protoAdminUserToDomainAdminUser(oResp.GetAdminUser())

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error) {

	oRequest := &pbResourceModel.AdminUserShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminUsers := make([]*domain.AdminUser, 0, len(oResponse.GetAdminUsers()))
	for _, oOne := range oResponse.GetAdminUsers() {
		oAdminUser := protoAdminUserToDomainAdminUser(oOne)
		aAdminUsers = append(aAdminUsers, &oAdminUser)
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminUserTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.TotalByFilters(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {

	oRequest := &pbResourceModel.AdminUserAddOneInput{
		Variable: &pbResourceModel.AdminUserVariable{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.AddOne(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	bStatus := oResponse.GetStatus()

	return bStatus, nil
}

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.AdminUserEditOneByIdInput{
		Id: uint32(iId),
		Variable: &pbResourceModel.AdminUserVariable{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.EditOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	bStatus := oResponse.GetStatus()

	return bStatus, nil
}

func (oSelf *AdminUserModel) RemoveOneById(iId uint) (bool, error) {

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminUserRemoveOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return false, oErr
	}

	bStatus := oResponse.GetStatus()

	return bStatus, nil
}
