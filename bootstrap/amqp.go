package bootstrap

import (
	"fmt"

	"github.com/charmbracelet/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewAmqp() (*amqp.Connection, error) {
	sDSN := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		CONFIG.AMQP.USERNAME,
		CONFIG.AMQP.PASSWORD,
		CONFIG.AMQP.HOST,
		CONFIG.AMQP.PORT,
	)

	oConnection, oErr := amqp.Dial(sDSN)

	log.Info("[INFO] AMQP 連線完成.", "addr", CONFIG.AMQP.HOST+":"+CONFIG.AMQP.PORT)

	return oConnection, oErr
}
