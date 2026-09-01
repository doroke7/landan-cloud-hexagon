package outputApplicationZeromqEvent

import (
	"encoding/json"

	"github.com/go-zeromq/zmq4"

	domain "example/internal/domain"
	outputApplicationZeromq "example/internal/output/application/zeromq"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationZeromq.AbstractZeromq
}

func NewAdminUserEvent(oAbstractZeromq *outputApplicationZeromq.AbstractZeromq) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractZeromq: oAbstractZeromq,
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
