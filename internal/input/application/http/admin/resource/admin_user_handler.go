package applicationHttpAdminResource

import (
	"github.com/gin-gonic/gin"

	pkg "example/pkg"

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

	oRequest := &pkg.Request{Context: oContext}

	oValue := &domain.AdminUserValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "request 格式錯誤", struct{}{}, 0, "")
		return
	}

	bResult, oErr := oSelf.AdminUserUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {
			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), "新增失敗", struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "新增失敗", struct{}{}, 0, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "新增成功", struct{}{}, 0, "")

}

func (oSelf *AdminUserHandler) ShowOnes(oContext *gin.Context) {

	oRequest := &pkg.Request{Context: oContext}

	oPagination := &pkg.Pagination{}
	if oErr := oRequest.Bind("pagination", oPagination); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "pagination 格式錯誤", struct{}{}, 0, "")
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

	var aFilters []*pkg.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "")
		return
	}

	var aSorters []*pkg.Sorter
	if oErr := oRequest.Bind("sorters", &aSorters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "sorters 格式錯誤", struct{}{}, 0, "")
		return
	}

	aAdminUsers, iTotal, oErr := oSelf.AdminUserUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Message, struct{}{}, 0, "")
			return
		}
		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, 0, "")
		return
	}

	oResult := pkg.NewResultOnes(aAdminUsers)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "")

}
