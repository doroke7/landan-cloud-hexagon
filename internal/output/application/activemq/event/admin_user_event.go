package outputApplicationActivemqEvent

import (
	"encoding/json"

	domain "example/internal/domain"
	outputApplicationActivemq "example/internal/output/application/activemq"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationActivemq.AbstractActivemq
}

func NewAdminUserEvent(oAbstractActivemq *outputApplicationActivemq.AbstractActivemq) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractActivemq: oAbstractActivemq,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	err = oSelf.Conn.Send("/Queue/AdminUser.AddOne", "application/json", aByteBody)

	return err
}

// Close 是空實作：Conn 現在是從 AbstractActivemq 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserEvent) Close() error {
	return nil
}
