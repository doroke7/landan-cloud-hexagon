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
	GameUsecase usecasePortAnyAdminResource.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oGameUsecase usecasePortAnyAdminResource.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler: oAbstractHandler,
		GameUsecase:     oGameUsecase,
	}
}

func (oSelf *GameHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &domain.GameValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "request 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	bResult, oErr := oSelf.GameUsecase.AddOne(oValue)
	_ = bResult

	if oErr != nil {
		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Error(), struct{}{}, 0, "", nil)
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "新增失敗", struct{}{}, 0, "", nil)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "新增成功", struct{}{}, 0, "", nil)

}

func (oSelf *GameHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", nil)
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

	oValue := &domain.GameValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "value 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.GameUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Error(), struct{}{}, 0, "", nil)
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "修改失敗", struct{}{}, 0, "", nil)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "修改成功", struct{}{}, 0, "", nil)

}

func (oSelf *GameHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", nil)
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
	bResult, oErr := oSelf.GameUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Error(), struct{}{}, 0, "", nil)
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "刪除失敗", struct{}{}, 0, "", nil)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "刪除成功", struct{}{}, 0, "", nil)

}

func (oSelf *GameHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 位置錯誤", struct{}{}, 0, "", nil)
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	iId := uint(fId)
	oGame, oErr := oSelf.GameUsecase.ShowOne(iId)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Error(), struct{}{}, 0, "", nil)
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "查詢失敗", struct{}{}, 0, "", nil)
		return
	}

	oResult := pkgGin.NewResult(oGame, nil, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, 0, "", nil)

}

func (oSelf *GameHandler) ShowOnes(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oPagination := &pkgInput.Pagination{}
	if oErr := oRequest.Bind("pagination", oPagination); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "pagination 格式錯誤", struct{}{}, 0, "", nil)
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
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	var aSorters []*pkgInput.Sorter
	if oErr := oRequest.Bind("sorters", &aSorters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "sorters 格式錯誤", struct{}{}, 0, "", nil)
		return
	}

	aGames, iTotal, oErr := oSelf.GameUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Error(), struct{}{}, 0, "", nil)
			return
		}
		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, 0, "", nil)
		return
	}

	oResult := pkgGin.NewResult(nil, aGames, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "", nil)

}
