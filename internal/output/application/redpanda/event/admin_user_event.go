package outputApplicationRedpandaEvent

import (
	"encoding/json"

	"github.com/IBM/sarama"

	domain "example/internal/domain"
	outputApplicationRedpanda "example/internal/output/application/redpanda"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*outputApplicationRedpanda.AbstractRedpanda
}

func NewAdminUserEvent(oAbstractRedpanda *outputApplicationRedpanda.AbstractRedpanda) (outputPortAnyEvent.AdminUserEvent, error) {
	return &AdminUserEvent{
		AbstractRedpanda: oAbstractRedpanda,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, _, err = oSelf.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "/Queue/AdminUser.AddOne",
		Value: sarama.ByteEncoder(aByteBody),
	})

	return err
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
