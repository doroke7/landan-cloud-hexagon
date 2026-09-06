package outputApplicationPulsarEvent

import (
	"encoding/json"

	"github.com/apache/pulsar-client-go/pulsar"

	domain "example/internal/domain"
	outputPortAnyEvent "example/internal/output/port/any/event"
)

type AdminUserEvent struct {
	*AbstractEvent
	Producer pulsar.Producer
}

func NewAdminUserEvent(oAbstractEvent *AbstractEvent) (outputPortAnyEvent.AdminUserEvent, error) {
	oProducer, err := oAbstractEvent.Client.CreateProducer(pulsar.ProducerOptions{
		Topic: "/Queue/AdminUser.AddOne",
	})
	if err != nil {
		return nil, err
	}

	return &AdminUserEvent{
		AbstractEvent: oAbstractEvent,
		Producer:      oProducer,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, err = oSelf.Producer.Send(oSelf.Context, &pulsar.ProducerMessage{
		Payload: aByteBody,
	})

	return err
}

func (oSelf *AdminUserEvent) Close() error {
	oSelf.Producer.Close()
	return nil
}
