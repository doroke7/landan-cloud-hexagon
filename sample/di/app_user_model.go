package main

type AppUserModel struct {
}

// 錯誤的 DI 範例
func NewAppUserModel() *AppUserModel {

	return &AppUserModel{}
}
