package inputApplicationHttpAdminOption

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"

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

func (oSelf *GameTypeHandler) SelectOnes(oContext *gin.Context) {

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

	aGameTypes, iTotal, oErr := oSelf.GameTypeUsecase.ShowOnes(aFilters, aSorters, oPagination)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, aGameTypes, nil)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, int(iTotal), "")

}

func (oSelf *GameTypeHandler) SelectTree(oContext *gin.Context) {

	aTree, oErr := oSelf.GameTypeUsecase.ShowTree()

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oResult := pkgGin.NewResult(nil, nil, aTree)

	oSelf.Response.Set(oContext, 200, 1, "成功查詢", oResult, len(aTree), "")

}
