package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	AdminPermissionGroupModel outputPortAnyModel.AdminPermissionGroupModel
	AdminPermissionGroupLogic outputPortAnyLogic.AdminPermissionGroupLogic
}

func NewAdminPermissionGroupUsecase(oAdminPermissionGroupModel outputPortAnyModel.AdminPermissionGroupModel, oAdminPermissionGroupLogic outputPortAnyLogic.AdminPermissionGroupLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.AdminPermissionGroupUsecase {
	return &AdminPermissionGroupUsecase{
		AbstractUsecase:           oAbstractUsecase,
		AdminPermissionGroupModel: oAdminPermissionGroupModel,
		AdminPermissionGroupLogic: oAdminPermissionGroupLogic,
	}
}

func (oSelf *AdminPermissionGroupUsecase) AddOne(oValue *domain.AdminPermissionGroupValue) error {

	oErr := oSelf.AdminPermissionGroupModel.AddOne(oValue)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) EditOne(oValue *domain.AdminPermissionGroupValue, iId uint64) error {

	oErr := oSelf.AdminPermissionGroupModel.EditOneById(oValue, iId)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) RemoveOne(iId uint64) error {

	// TY oSelf.AdminPermissionGroupModel

	oErr := oSelf.AdminPermissionGroupModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowOne(iId uint64) (*domain.AdminPermissionGroup, error) {

	oAdminPermissionGroup, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroupById(iId)

	return oAdminPermissionGroup, oErr
}

func (oSelf *AdminPermissionGroupUsecase) ShowTree() ([]*domain.AdminPermissionGroup, error) {

	aAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupLogic.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {

	aAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroups()

	iTotal := len(aAdminPermissionGroups)

	return aAdminPermissionGroups, uint64(iTotal), oErr
}
