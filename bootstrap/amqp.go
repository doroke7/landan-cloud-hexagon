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

	oConnection, oErr := amqp.Dial(sDSN)

	fmt.Println("[INFO] AMQP 連線完成. ", CONFIG.AMQP.HOST+":"+CONFIG.AMQP.PORT)

	return oConnection, oErr
}
