package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
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

func (oSelf *AdminPermissionGroupUsecase) AddOne(oVariable *domain.AdminPermissionGroupVariable) error {

	if oVariable.ParentId != nil && *oVariable.ParentId != 0 {
		oParentAdminPermissionGroup, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroupById(*oVariable.ParentId)
		if oErr != nil {
			return oErr
		}

		if oParentAdminPermissionGroup == nil {
			oParentNotFoundError := pkgUtility.NewDefaultError("parent admin permission group not found", -2, 200)

			return oParentNotFoundError
		}

		if oParentAdminPermissionGroup.ParentId != 0 {
			oExceedDepthError := pkgUtility.NewDefaultError("admin permission group depth cannot exceed 2 levels", -2, 200)

			return oExceedDepthError
		}
	}

	oErr := oSelf.AdminPermissionGroupLogic.AddAdminPermissionGroup(oVariable)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) EditOne(oVariable *domain.AdminPermissionGroupVariable, iId uint64) error {

	if oVariable.ParentId != nil && *oVariable.ParentId != 0 {
		if *oVariable.ParentId == iId {
			oSelfParentError := pkgUtility.NewDefaultError("admin permission group cannot be its own parent", -2, 200)

			return oSelfParentError
		}

		oParentAdminPermissionGroup, oErr := oSelf.AdminPermissionGroupLogic.ShowAdminPermissionGroupById(*oVariable.ParentId)
		if oErr != nil {
			return oErr
		}

		if oParentAdminPermissionGroup == nil {
			oParentNotFoundError := pkgUtility.NewDefaultError("parent admin permission group not found", -2, 200)

			return oParentNotFoundError
		}

		if oParentAdminPermissionGroup.ParentId != 0 {
			oExceedDepthError := pkgUtility.NewDefaultError("admin permission group depth cannot exceed 2 levels", -2, 200)

			return oExceedDepthError
		}

		aChildAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupModel.ShowOnesByParentId(iId)
		if oErr != nil {
			return oErr
		}

		if len(aChildAdminPermissionGroups) >= 1 {
			oHasChildrenError := pkgUtility.NewDefaultError("admin permission group has child groups; moving it would exceed 2 levels", -2, 200)

			return oHasChildrenError
		}
	}

	oErr := oSelf.AdminPermissionGroupLogic.EditAdminPermissionGroupById(oVariable, iId)

	return oErr
}

func (oSelf *AdminPermissionGroupUsecase) RemoveOne(iId uint64) error {

	aAdminPermissionGroups, oErr := oSelf.AdminPermissionGroupModel.ShowOnesByParentId(iId)

	if oErr != nil {
		return oErr
	}

	if len(aAdminPermissionGroups) >= 1 {
		oHasChildrenError := pkgUtility.NewDefaultError("Cannot delete because this group still has child groups", -2, 200)

		return oHasChildrenError
	}

	oErr = oSelf.AdminPermissionGroupModel.RemoveOneById(iId)

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
