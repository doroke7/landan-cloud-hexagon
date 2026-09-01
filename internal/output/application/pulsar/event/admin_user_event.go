package outputApplicationPulsarEvent

import (
	"encoding/json"

	"github.com/apache/pulsar-client-go/pulsar"

	domain "example/internal/domain"
	outputApplicationPulsar "example/internal/output/application/pulsar"
)

type AdminUserEvent struct {
	*outputApplicationPulsar.AbstractPulsar
	Producer pulsar.Producer
}

func NewAdminUserEvent(oAbstractPulsar *outputApplicationPulsar.AbstractPulsar) (*AdminUserEvent, error) {
	oProducer, err := oAbstractPulsar.Client.CreateProducer(pulsar.ProducerOptions{
		Topic: "/Queue/AdminUser.AddOne",
	})
	if err != nil {
		return nil, err
	}

	return &AdminUserEvent{
		AbstractPulsar: oAbstractPulsar,
		Producer:       oProducer,
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
