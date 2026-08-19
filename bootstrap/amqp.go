package bootstrap

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewAmqp() (*amqp.Connection, error) {
	sDSN := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		CONFIG.AMQP.USER,
		CONFIG.AMQP.PASS,
		CONFIG.AMQP.HOST,
		CONFIG.AMQP.PORT,
	)

	return amqp.Dial(sDSN)
}
