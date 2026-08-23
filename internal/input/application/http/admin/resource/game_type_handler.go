package inputApplicationHttpAdminResource

import (
	"github.com/gin-gonic/gin"

	pkg "example/pkg"

	domain "example/internal/domain"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
)

type GameTypeHandler struct {
	*inputApplicationHttp.AbstractHandler
	GameTypeUsecase usecasePortAnyAdminResource.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oGameTypeUsecase usecasePortAnyAdminResource.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler: oAbstractHandler,
		GameTypeUsecase: oGameTypeUsecase,
	}
}

func (oSelf *GameTypeHandler) AddOne(oContext *gin.Context) {

	oRequest := &pkg.Request{Context: oContext}

	oValue := &domain.GameTypeValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "request 格式錯誤", struct{}{}, 0, "")
		return
	}

	bResult, oErr := oSelf.GameTypeUsecase.AddOne(oValue)
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

func (oSelf *GameTypeHandler) EditOne(oContext *gin.Context) {

	oRequest := &pkg.Request{Context: oContext}

	var aFilters []*pkg.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "")
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter 位置錯誤", struct{}{}, 0, "")
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "")
		return
	}

	oValue := &domain.GameTypeValue{}
	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "value 格式錯誤", struct{}{}, 0, "")
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.GameTypeUsecase.EditOne(oValue, iId)
	_ = bResult

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), "修改失敗", struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "修改失敗", struct{}{}, 0, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "修改成功", struct{}{}, 0, "")

}

func (oSelf *GameTypeHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkg.Request{Context: oContext}

	var aFilters []*pkg.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "")
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter 位置錯誤", struct{}{}, 0, "")
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "")
		return
	}

	iId := uint(fId)
	bResult, oErr := oSelf.GameTypeUsecase.RemoveOne(iId)
	_ = bResult

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), "刪除失敗", struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "刪除失敗", struct{}{}, 0, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "刪除成功", struct{}{}, 0, "")

}

func (oSelf *GameTypeHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkg.Request{Context: oContext}

	var aFilters []*pkg.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "filters 格式錯誤", struct{}{}, 0, "")
		return
	}

	if aFilters == nil || len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 位置錯誤", struct{}{}, 0, "")
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		oSelf.Response.Set(oContext, 200, -1, "filter.id 格式錯誤", struct{}{}, 0, "")
		return
	}

	iId := uint(fId)
	oGameType, oErr := oSelf.GameTypeUsecase.ShowOne(iId)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), "查詢失敗", struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, "查詢失敗", struct{}{}, 0, "")
		return
	}

	oResult := pkg.NewResultOne(oGameType)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, 0, "")

}

func (oSelf *GameTypeHandler) ShowOnes(oContext *gin.Context) {

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

	aGameTypes, iTotal, oErr := oSelf.GameTypeUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {

			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Message, struct{}{}, 0, "")
			return
		}
		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, 0, "")
		return
	}

	oResult := pkg.NewResultOnes(aGameTypes)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "")

}
