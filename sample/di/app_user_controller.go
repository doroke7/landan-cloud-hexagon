package main

// 關於DI：要注意的事情是
// 上面 struct 結構體只是定義了結構
// 真正注入的地方在 下面的 NewStruct 函數

type AppUserController struct {
	*AppUserModel
	*AppUserLogic
}

// 錯誤的 DI 範例
// 錯誤寫法：把 NewAppUserModel （注入物件）寫在 NewMovieController（被注入物件） 裡面
func NewAppUserController() *AppUserController {
	return &AppUserController{
		AppUserModel: NewAppUserModel(),
		AppUserLogic: NewAppUserLogic(),
	}
}
