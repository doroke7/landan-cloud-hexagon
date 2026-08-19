package main

import (
	"go.uber.org/dig"
)

func Container() *dig.Container {
	oContainer := dig.New()

	// 1. 把所有元件的生產說明書全部「鋪平」註冊進去
	_ = oContainer.Provide(NewMovieModel)
	_ = oContainer.Provide(NewMovieLogic)
	_ = oContainer.Provide(NewMovieController)

	return oContainer

}
