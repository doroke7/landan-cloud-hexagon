package outputApplicationZeromqEvent

import (
	"encoding/json"

	"github.com/go-zeromq/zmq4"

	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*AbstractEvent
}

func NewAdminUserEvent(oAbstractEvent *AbstractEvent) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractEvent: oAbstractEvent,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Socket.SendMulti(zmq4.NewMsgFrom(
		[]byte("/Queue/AdminUser.AddOne"),
		aByteBody,
	))
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
