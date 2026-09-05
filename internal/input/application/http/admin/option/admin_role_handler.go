package inputApplicationHttpAdminOption

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
)

type AdminRoleHandler struct {
	*inputApplicationHttp.AbstractHandler
	AdminOptionAdminRoleUsecase usecasePortAnyAdminOption.AdminRoleUsecase
}

func NewAdminRoleHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAdminRoleUsecase usecasePortAnyAdminOption.AdminRoleUsecase) *AdminRoleHandler {
	return &AdminRoleHandler{
		AbstractHandler:             oAbstractHandler,
		AdminOptionAdminRoleUsecase: oAdminRoleUsecase,
	}
}

func (oSelf *AdminRoleHandler) SelectOnes(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oPagination := &pkgInput.Pagination{}
	if oErr := oRequest.Bind("pagination", oPagination); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("pagination format error", -1, 200))
		return
	}
	if oPagination.Size == nil || *oPagination.Size == 0 {
		iSize := uint(10)
		oPagination.Size = &iSize
	}
	if oPagination.Page == nil || *oPagination.Page == 0 {
		iPage := uint(1)
		oPagination.Page = &iPage
	}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	var aSorters []*pkgInput.Sorter
	if oErr := oRequest.Bind("sorters", &aSorters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("sorters format error", -1, 200))
		return
	}

	aAdminRoles, iTotal, oErr := oSelf.AdminOptionAdminRoleUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aAdminRoles, nil)

	oSelf.Response.Set(oContext, 200, 1, "Query successful", oResult, int(iTotal), "")

}
