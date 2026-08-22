package rabbitmq

import (
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"

	domain "example/internal/domain"
	outputApplicationRabbitmq "example/internal/output/application/rabbitmq"
)

type AdminUserModel struct {
	*outputApplicationRabbitmq.AbstractRabbitmq
}

func NewAdminUserModel(oAbstractRabbitmq *outputApplicationRabbitmq.AbstractRabbitmq) (*AdminUserModel, error) {
	_, err := oAbstractRabbitmq.Channel.QueueDeclare("AdminUser.AddOne", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &AdminUserModel{
		AbstractRabbitmq: oAbstractRabbitmq,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
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

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by rabbitmq")
}

// Close 是空實作：Conn／Channel 現在都是從 AbstractRepository 注入的共用資源，
// 生命週期不屬於這個 repository，不該由這裡關閉。
func (oSelf *AdminUserModel) Close() error {
	return nil
}
