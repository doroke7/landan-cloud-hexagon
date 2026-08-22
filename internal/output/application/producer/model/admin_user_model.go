package producer

import (
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"

	domain "example/internal/domain"
	producerBase "example/internal/output/application/producer"
)

type AdminUserModel struct {
	*producerBase.AbstractProducer
}

func NewAdminUserModel(oAbstractProducer *producerBase.AbstractProducer) (*AdminUserModel, error) {
	_, err := oAbstractProducer.Channel.QueueDeclare("AdminUser.AddOne", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &AdminUserModel{
		AbstractProducer: oAbstractProducer,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	body, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	return oSelf.Channel.Publish(
		"",
		"AdminUser.AddOne",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by producer")
}

// Close 是空實作：Conn／Channel 現在都是從 AbstractRepository 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
