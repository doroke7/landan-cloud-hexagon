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

	oValue := &domain.GameVariable{}
	if oErr := oRequest.Bind("variable", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("request format error", -1, 200))
		return
	}

	oErr := oSelf.AdminResourceGameUsecase.AddOne(oValue)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Added successfully", struct{}{}, 0, "")

}

func (oSelf *GameHandler) EditOne(oContext *gin.Context) {

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

	oValue := &domain.GameVariable{}
	if oErr := oRequest.Bind("variable", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("value format error", -1, 200))
		return
	}

	iId := uint64(fId)
	oErr := oSelf.AdminResourceGameUsecase.EditOne(oValue, iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Updated successfully", struct{}{}, 0, "")

}

func (oSelf *GameHandler) RemoveOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	if aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter position error", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id format error", -1, 200))
		return
	}

	iId := uint(fId)
	oErr := oSelf.AdminResourceGameUsecase.RemoveOne(iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Deleted successfully", struct{}{}, 0, "")

}

func (oSelf *GameHandler) ShowOne(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	var aFilters []*pkgInput.Filter
	if oErr := oRequest.Bind("filters", &aFilters); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("filters format error", -1, 200))
		return
	}

	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].Field == nil || *aFilters[0].Field != "id" {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id position error", -1, 200))
		return
	}

	fId, bOk := aFilters[0].Value.(float64)
	if !bOk {
		_ = oContext.Error(pkgUtility.NewDefaultError("filter.id format error", -1, 200))
		return
	}

	iId := uint64(fId)
	oGame, oErr := oSelf.AdminResourceGameUsecase.ShowOne(iId)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(oGame, nil, nil)

	oSelf.Response.Set(oContext, 200, 1, "Query successful", oResult, 0, "")

}

func (oSelf *GameHandler) ShowOnes(oContext *gin.Context) {

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

	aGames, iTotal, oErr := oSelf.AdminResourceGameUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aGames, nil)

	oSelf.Response.Set(oContext, 200, 1, "Query successful", oResult, int(iTotal), "")

}
