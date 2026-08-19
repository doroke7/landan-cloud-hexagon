package main

// 關於DI：要注意的事情是
// 上面 struct 結構體只是定義了結構
// 真正注入的地方在 下面的 NewStruct 函數

type MovieController struct {
	*MovieModel
	*MovieLogic
}

// 錯誤寫法：把 NewAppUserModel （注入物件）寫在 NewMovieController（被注入物件） 裡面
func NewMovieController(oMovieModel *MovieModel, oMovieLogic *MovieLogic) *MovieController {
	return &MovieController{
		MovieModel: oMovieModel,
		MovieLogic: oMovieLogic,
	}
}
