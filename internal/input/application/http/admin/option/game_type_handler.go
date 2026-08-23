package inputApplicationHttpAdminOption

import (
	"github.com/gin-gonic/gin"

	pkg "example/pkg"

	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminOption "example/internal/usecase/port/any/admin/option"
)

type GameTypeHandler struct {
	*inputApplicationHttp.AbstractHandler
	GameTypeUsecase usecasePortAnyAdminOption.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oGameTypeUsecase usecasePortAnyAdminOption.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler: oAbstractHandler,
		GameTypeUsecase: oGameTypeUsecase,
	}
}

func (oSelf *GameTypeHandler) Select(oContext *gin.Context) {

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
