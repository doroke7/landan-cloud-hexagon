package main

type AppUserLogic struct {
	*AppUserModel
}

// 錯誤的 DI 範例
func NewAppUserLogic() *AppUserLogic {

	return &AppUserLogic{
		AppUserModel: NewAppUserModel(),
	}
}
