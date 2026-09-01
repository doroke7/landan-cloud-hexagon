package usecaseApplicationAnyEvent

import (
	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
	usecasePortAnyEvent "example/internal/usecase/port/any/event"
)

type AdminUserUsecase struct {
	outputPortAnyEvent.AdminUserEvent
}

func NewAdminUserUsecase(oAdminUserEvent outputPortAnyEvent.AdminUserEvent) usecasePortAnyEvent.AdminUserUsecase {
	return &AdminUserUsecase{
		AdminUserEvent: oAdminUserEvent,
	}
}

func (oSelf *AdminUserUsecase) AddOne(oAdminUser *domain.AdminUser) error {
	return oSelf.AdminUserEvent.AddOne(oAdminUser)
}
