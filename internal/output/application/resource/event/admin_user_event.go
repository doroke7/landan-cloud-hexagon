package outputApplicationResourceEvent

import (
	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
	pb "example/pb"
	pbResourceEvent "example/pb/resource/event"
)

type AdminUserEvent struct {
	*AbstractEvent
}

func NewAdminUserEvent(oAbstractEvent *AbstractEvent) (outputPortAnyEvent.AdminUserEvent, error) {
	oEvent := &AdminUserEvent{
		AbstractEvent: oAbstractEvent,
	}

	return oEvent, nil
}

// AddOne 打 resource gRPC 的 AdminUserEvent.AddOne；event port 只回 error。
func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	iId := oAdminUser.Id
	sName := oAdminUser.Name
	sPassword := oAdminUser.Password

	oValue := &pb.AppUserVariable{
		Id:       &iId,
		Name:     &sName,
		Password: &sPassword,
	}
	oRequest := &pbResourceEvent.AdminUserEventAddOneInput{
		Variable: oValue,
	}
	_, oErr := oSelf.ResourceEventClient.AdminUser.AddOne(oSelf.Context, oRequest)

	return oErr
}
