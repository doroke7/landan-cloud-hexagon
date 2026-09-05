package inputApplicationHttpAdminResource

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

	domain "example/internal/domain"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AdminUserHandler struct {
	*inputApplicationHttp.AbstractHandler
	AdminResourceAdminUserUsecase usecasePortAnyAdminResource.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAdminUserUsecase usecasePortAnyAdminResource.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:               oAbstractHandler,
		AdminResourceAdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &domain.AdminUserValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("request format error", -1, 200))
		return
	}

	bResult, oErr := oSelf.AdminResourceAdminUserUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Added successfully", struct{}{}, 0, "")

}

func (oSelf *AdminUserHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter position error", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id format error", -1, 200))
		return
	}

	oValue := &domain.AdminUserValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("value format error", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminResourceAdminUserUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Updated successfully", struct{}{}, 0, "")

}

func (oSelf *AdminUserHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter position error", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id format error", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminResourceAdminUserUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Deleted successfully", struct{}{}, 0, "")

}

func (oSelf *AdminUserHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter position error", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id format error", -1, 200))
		return
	}

	iId := uint(fId)
	oAdminUser, oErr := oSelf.AdminResourceAdminUserUsecase.ShowOne(iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(oAdminUser, nil, nil)

	oSelf.Response.Set(oContext, 200, 1, "Query successful", oResult, 0, "")

}

func (oSelf *AdminUserHandler) ShowOnes(oContext *gin.Context) {

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

	aAdminUsers, iTotal, oErr := oSelf.AdminResourceAdminUserUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	for _, oAdminUser := range aAdminUsers {
		if oAdminUser != nil {
			oAdminUser.Password = "" // 列表不外洩密碼
		}
	}

	oResult := pkgGin.NewResult(nil, aAdminUsers, nil)

	oSelf.Response.Set(oContext, 200, 1, "Query successful", oResult, int(iTotal), "")

}
