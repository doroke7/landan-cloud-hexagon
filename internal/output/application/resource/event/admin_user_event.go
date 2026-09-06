package outputApplicationResourceEvent

import (
	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
	pbResourceEvent "example/pb/resource/event"
)

type AdminUserEvent struct {
	*AbstractEvent
}

func NewAdminUserEvent(oAbstractEvent *AbstractEvent) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractEvent: oAbstractEvent,
	}, nil
}

// AddOne 打 resource gRPC 的 AdminUserEvent.AddOne；event port 只回 error。
func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	iId := oAdminUser.Id
	sName := oAdminUser.Name
	sPassword := oAdminUser.Password

	_, oErr := oSelf.ResourceEventClient.AdminUser.AddOne(oSelf.Context, &pbResourceEvent.AdminUserEventAddOneInput{
		Value: &pbResourceEvent.AppUserValue{
			Id:       &iId,
			Name:     &sName,
			Password: &sPassword,
		},
	})

	return oErr
}
