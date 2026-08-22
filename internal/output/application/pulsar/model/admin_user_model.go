package pulsar

import (
	"encoding/json"
	"errors"

	"github.com/apache/pulsar-client-go/pulsar"

	domain "example/internal/domain"
	outputApplicationPulsar "example/internal/output/application/pulsar"
)

type AdminUserModel struct {
	*outputApplicationPulsar.AbstractPulsar
	Producer pulsar.Producer
}

func NewAdminUserModel(oAbstractPulsar *outputApplicationPulsar.AbstractPulsar) (*AdminUserModel, error) {
	oProducer, err := oAbstractPulsar.Client.CreateProducer(pulsar.ProducerOptions{
		Topic: "/Queue/AdminUser.AddOne",
	})
	if err != nil {
		return nil, err
	}

	return &AdminUserModel{
		AbstractPulsar: oAbstractPulsar,
		Producer:       oProducer,
	}, nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUser) error {
	aByteBody, err := json.Marshal(oAdminUser)
	if err != nil {
		return err
	}

	_, err = oSelf.Producer.Send(oSelf.Context, &pulsar.ProducerMessage{
		Payload: aByteBody,
	})

	return err
}

func (oSelf *AdminUserModel) ShowOneById(id uint) (*domain.AdminUser, error) {
	return nil, errors.New("not supported by pulsar")
}

// Close 這裡不是空實作：Producer 是這個 model 自己在建構時針對
// "AdminUser.AddOne" topic 開的，跟 rabbitmq/redpanda 那種掛在
// AbstractXxx 上的共用資源不同，沒有其他人共用，生命週期歸這裡管。
func (oSelf *AdminUserModel) Close() error {
	oSelf.Producer.Close()
	return nil
}
