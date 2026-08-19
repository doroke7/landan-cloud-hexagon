package main

type MovieLogic struct {
	*MovieModel
}

func NewMovieLogic(oMovieModel *MovieModel) *MovieLogic {
	return &MovieLogic{
		MovieModel: oMovieModel,
	}
}
