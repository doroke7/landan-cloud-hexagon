package inputApplicationHttpAdminResource

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgInput "example/pkg/input"

	domain "example/internal/domain"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type AdminUserHandler struct {
	*inputApplicationHttp.AbstractHandler
	AdminUserUsecase usecasePortAnyAdminResource.AdminUserUsecase
}

func NewAdminUserHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAdminUserUsecase usecasePortAnyAdminResource.AdminUserUsecase) *AdminUserHandler {
	return &AdminUserHandler{
		AbstractHandler:  oAbstractHandler,
		AdminUserUsecase: oAdminUserUsecase,
	}
}

func (oSelf *AdminUserHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &domain.AdminUserValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "request 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	bResult, oErr := oSelf.AdminUserUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "新增成功", struct{}{}, 0, "", nil)

}

func (oSelf *AdminUserHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter 位置錯誤", struct{}{}, 0, "", nil)
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	oValue := &domain.AdminUserValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "value 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminUserUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "修改成功", struct{}{}, 0, "", nil)

}

func (oSelf *AdminUserHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter 位置錯誤", struct{}{}, 0, "", nil)
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminUserUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "刪除成功", struct{}{}, 0, "", nil)

}

func (oSelf *AdminUserHandler) ShowOnes(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oPagination := &pkgInput.Pagination{}
	if oErr := oRequest.Bind("pagination", oPagination); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "pagination 格式錯誤", struct{}{}, 0, "", oErr)
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
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	var aSorters []*pkgInput.Sorter
	if oErr := oRequest.Bind("sorters", &aSorters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "sorters 格式錯誤", struct{}{}, 0, "", oErr)
		return
	}

	aAdminUsers, iTotal, oErr := oSelf.AdminUserUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aAdminUsers, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "", nil)

}
