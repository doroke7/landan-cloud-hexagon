package outputApplicationResourceModel

import (
	"errors"

	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

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

	oProtoAdminUser := oResp.GetAdminUser()

	return &domain.AdminUser{
		Id:       uint(oProtoAdminUser.GetId()),
		Name:     oProtoAdminUser.GetName(),
		Password: oProtoAdminUser.GetPassword(),
	}, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {

	oResp, err := oSelf.ResourceModelClient.AdminUser.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminUserShowOneByIdInput{Id: uint32(iId)},
	)

	if err != nil {
		return nil, err
	}

	oProtoAdminUser := oResp.GetAdminUser()

	return &domain.AdminUser{
		Id:       uint(oProtoAdminUser.GetId()),
		Name:     oProtoAdminUser.GetName(),
		Password: oProtoAdminUser.GetPassword(),
	}, nil
}

// ShowOnesByFiltersWithSortersPagination 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error) {
	return nil, errors.New("not supported by resource 1")
}

// TotalByFilters 目前 Resource gRPC service 沒有對應的 RPC，先不支援。
func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	return 0, errors.New("not supported by resource 2")
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {

	oRequest := &pbResourceModel.AdminUserAddOneInput{
		Varible: &pbResourceModel.AdminUserVarible{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.AddOne(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.AdminUserEditOneByIdInput{
		Id: uint32(iId),
		Varible: &pbResourceModel.AdminUserVarible{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.EditOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}
