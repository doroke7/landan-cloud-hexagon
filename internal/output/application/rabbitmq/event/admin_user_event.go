package outputApplicationRabbitmqEvent

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	domain "example/internal/domain"
	outputApplicationRabbitmq "example/internal/output/application/rabbitmq"
)

type AdminUserEvent struct {
	*outputApplicationRabbitmq.AbstractRabbitmq
}

func NewAdminUserEvent(oAbstractRabbitmq *outputApplicationRabbitmq.AbstractRabbitmq) (*AdminUserEvent, error) {
	_, err := oAbstractRabbitmq.Channel.QueueDeclare("AdminUser.AddOne", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &AdminUserEvent{
		AbstractRabbitmq: oAbstractRabbitmq,
	}, nil
}

func (oSelf *AdminUserEvent) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Channel.Publish(
		"",
		"/Queue/AdminUser.AddOne",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        aByteBody,
		},
	)
}

func (oSelf *AdminUserEvent) Close() error {
	return nil
}
