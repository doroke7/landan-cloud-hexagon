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

type GameHandler struct {
	*inputApplicationHttp.AbstractHandler
	AdminResourceGameUsecase usecasePortAnyAdminResource.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oGameUsecase usecasePortAnyAdminResource.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler:          oAbstractHandler,
		AdminResourceGameUsecase: oGameUsecase,
	}
}

func (oSelf *GameHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &domain.GameValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("request 格式錯誤", -1, 200))
		return
	}

	bResult, oErr := oSelf.AdminResourceGameUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "新增成功", struct{}{}, 0, "")

}

func (oSelf *GameHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	oValue := &domain.GameValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("value 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminResourceGameUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "修改成功", struct{}{}, 0, "")

}

func (oSelf *GameHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.AdminResourceGameUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "刪除成功", struct{}{}, 0, "")

}

func (oSelf *GameHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters 格式錯誤", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 位置錯誤", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id 格式錯誤", -1, 200))
		return
	}

	iId := uint(fId)
	oGame, oErr := oSelf.AdminResourceGameUsecase.ShowOne(iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(oGame, nil, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, 0, "")

}

func (oSelf *GameHandler) ShowOnes(oContext *gin.Context) {

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

	aGames, iTotal, oErr := oSelf.AdminResourceGameUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aGames, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "")

}
