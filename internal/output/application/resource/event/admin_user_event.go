package outputApplicationResourceEvent

import (
	domain "example/internal/domain"
	outputApplicationResource "example/internal/output/application/resource"
	outputPortAnyEvent "example/internal/output/port/any/event"
	pbResourceEvent "example/pb/resource/event"
)

type AdminUserEvent struct {
	*outputApplicationResource.AbstractResource
}

func NewAdminUserEvent(oAbstractEvent *outputApplicationResource.AbstractResource) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractResource: oAbstractEvent,
	}, nil
}

// AddOne 打 resource gRPC 的 AdminUserEvent.AddOne；event port 只回 error。
func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	_, oErr := oSelf.ResourceEventClient.AdminUser.AddOne(oSelf.Context, &pbResourceEvent.AdminUserEventAddOneInput{
		Id:       uint64(oAdminUser.Id),
		Name:     oAdminUser.Name,
		Password: oAdminUser.Password,
	})

	return oErr
}
