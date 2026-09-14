package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pb "example/pb"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type AdminRoleModel struct {
	*AbstractModel
}

func NewAdminRoleModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminRoleModel {
	oModel := &AdminRoleModel{
		AbstractModel: oAbstractModel,
	}

	return oModel
}

func domainAdminRoleValueToProtoAdminRoleValue(oValue *domain.AdminRoleVariable) *pb.AdminRoleVariable {
	oAdminRoleValue := &pb.AdminRoleVariable{
		Key:  oValue.Key,
		Name: oValue.Name,
	}

	return oAdminRoleValue
}

func (oSelf *AdminRoleModel) AddOne(oAdminRole *domain.AdminRoleVariable) error {

	oRequest := &pbResourceModel.AdminRoleAddOneInput{
		Variable: domainAdminRoleValueToProtoAdminRoleValue(oAdminRole),
	}

	_, oErr := oSelf.ResourceModelClient.AdminRole.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminRoleModel) ShowOneById(iId uint64) (*domain.AdminRole, error) {

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.ShowOneById(
		oSelf.Context,
		&pbResourceModel.AdminRoleShowOneByIdInput{Id: uint64(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoAdminRole := oResponse.GetAdminRole()
	if oProtoAdminRole.GetId() == 0 {
		return nil, nil
	}

	oAdminRole := pkgProtoToDomain.AdminRole(oProtoAdminRole)

	return &oAdminRole, nil
}

func (oSelf *AdminRoleModel) EditOneById(oAdminRole *domain.AdminRoleVariable, iId uint64) error {

	oRequest := &pbResourceModel.AdminRoleEditOneByIdInput{
		Id:       uint64(iId),
		Variable: domainAdminRoleValueToProtoAdminRoleValue(oAdminRole),
	}

	_, oErr := oSelf.ResourceModelClient.AdminRole.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *AdminRoleModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.AdminRole.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.AdminRoleRemoveOneByIdInput{Id: uint64(iId)},
	)

	return oErr
}

func (oSelf *AdminRoleModel) ShowOnes() ([]*domain.AdminRole, error) {

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.ShowOnes(oSelf.Context, &pbResourceModel.AdminRoleShowOnesInput{})

	if oErr != nil {
		return nil, oErr
	}

	aAdminRoles := make([]*domain.AdminRole, 0, len(oResponse.GetAdminRoles()))
	for _, oOne := range oResponse.GetAdminRoles() {
		oAdminRole := pkgProtoToDomain.AdminRole(oOne)
		aAdminRoles = append(aAdminRoles, &oAdminRole)
	}

	return aAdminRoles, nil
}

func (oSelf *AdminRoleModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error) {

	oRequest := &pbResourceModel.AdminRoleShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aAdminRoles := make([]*domain.AdminRole, 0, len(oResponse.GetAdminRoles()))
	for _, oOne := range oResponse.GetAdminRoles() {
		oAdminRole := pkgProtoToDomain.AdminRole(oOne)
		aAdminRoles = append(aAdminRoles, &oAdminRole)
	}

	return aAdminRoles, nil
}

func (oSelf *AdminRoleModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.AdminRoleTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.AdminRole.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}
