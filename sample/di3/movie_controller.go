package main

// 關於DI：要注意的事情是
// 上面 struct 結構體只是定義了結構
// 真正注入的地方在 下面的 NewStruct 函數

type MovieController struct {
	*MovieModel
	*MovieLogic
}

func NewMovieController() *MovieController {
	return &MovieController{}
}
