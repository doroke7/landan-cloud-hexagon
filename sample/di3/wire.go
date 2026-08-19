//go:build wireinject

package main

import (
	"github.com/google/wire"
)

// 定義一個 Root Object, 便於取得所有 物件
type Container struct {
	*MovieModel
	*MovieLogic
	*MovieController
}

func InitContainer() *Container {
	wire.Build(
		NewMovieModel,
		NewMovieLogic,
		NewMovieController,
		wire.Struct(new(Container), "*"),
	)
	return nil
}
