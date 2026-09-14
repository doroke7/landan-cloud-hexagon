package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pb "example/pb"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminUserModel struct {
	*AbstractModel
}

func NewAdminUserModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractModel: oAbstractModel,
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

	oAdminUser := pkgProtoToDomain.AdminUser(oResp.GetAdminUser())

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint64) (*domain.AdminUser, error) {

	oResp, err := oSelf.ResourceModelClient.AdminUser.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminUserShowOneByIdInput{Id: uint32(iId)},
	)

	if err != nil {
		return nil, err
	}

	oAdminUser := pkgProtoToDomain.AdminUser(oResp.GetAdminUser())

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
		oAdminUser := pkgProtoToDomain.AdminUser(oOne)
		aAdminUsers = append(aAdminUsers, &oAdminUser)
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminUserTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminUser.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserVariable) error {

	oRequest := &pbResourceModel.AdminUserAddOneInput{
		Variable: &pb.AdminUserVariable{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminUser.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserVariable, iId uint64) error {

	oRequest := &pbResourceModel.AdminUserEditOneByIdInput{
		Id: iId,
		Variable: &pb.AdminUserVariable{
			Name:     oAdminUser.Name,
			Password: oAdminUser.Password,
		},
	}

	_, oErr := oSelf.ResourceModelClient.AdminUser.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminUserModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.AdminUser.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminUserRemoveOneByIdInput{Id: iId},
	)

	return oErr
}
