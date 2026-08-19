package main

import "fmt"

func main() {

	// 錯誤的 DI 範例， 每一個 AppUser 其實還是 不同 object
	oAppUserController := NewAppUserController()

	fmt.Printf("AppUser controller: %v\n", oAppUserController)

	///////////////////////////////////////////////////////////////////

	// 正確的 Di, 全部的 Movie 數據模型共用
	oMovieModel := NewMovieModel()
	oMovieController := NewMovieController(oMovieModel, NewMovieLogic(oMovieModel))
	fmt.Printf("Movie controller: %v\n", oMovieController)

	oContainer := Container()

	err := oContainer.Invoke(func(
		oMovieController *MovieController,
	) {
		// 在 Invoke 世界裡面就可以使用任何 有被註冊的 元件 object型態了
		// 而且 是已經幫你綁定 依賴注入 的嵌套規則了
		// 如 r.Get('/Admin/Resoure/Movie/showOnes', oMovieController.ShowOnes)
	})

	if err != nil {
		panic(fmt.Sprintf("❌ 容器跨套件注入或路由註冊失敗：%v", err))
	}

}
