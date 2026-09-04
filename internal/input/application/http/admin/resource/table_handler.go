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

type TableHandler struct {
	*inputApplicationHttp.AbstractHandler
	TableUsecase usecasePortAnyAdminResource.TableUsecase
}

func NewTableHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oTableUsecase usecasePortAnyAdminResource.TableUsecase) *TableHandler {
	return &TableHandler{
		AbstractHandler: oAbstractHandler,
		TableUsecase:    oTableUsecase,
	}
}

func (oSelf *TableHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &domain.TableValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("value 格式錯誤", -1, 200))
		return
	}

	bResult, oErr := oSelf.TableUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "新增成功", struct{}{}, 0, "", nil)

}

func (oSelf *TableHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	oTable, oErr := oSelf.TableUsecase.ShowOne(iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(oTable, nil, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, 0, "", nil)

}

func (oSelf *TableHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	oValue := &domain.TableValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("value 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.TableUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "修改成功", struct{}{}, 0, "", nil)

}

func (oSelf *TableHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.TableUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "刪除成功", struct{}{}, 0, "", nil)

}

func (oSelf *TableHandler) ShowOnes(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oPagination := &pkgInput.Pagination{}
	if oErr := oRequest.Bind("pagination", oPagination); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("pagination 格式錯誤", -1, 200))
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
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	var aSorters []*pkgInput.Sorter
	if oErr := oRequest.Bind("sorters", &aSorters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("sorters 格式錯誤", -1, 200))
		return
	}

	aTables, iTotal, oErr := oSelf.TableUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aTables, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "", nil)

}
