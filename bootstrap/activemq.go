package bootstrap

import (
	"fmt"
	"time"

	"github.com/go-stomp/stomp/v3"
)

// ActiveMQ 走 STOMP protocol 連線，跟 amqp.go 用 rabbitmq wire protocol 是不同協定。
func NewActivemq() (*stomp.Conn, error) {
	sAddr := fmt.Sprintf("%s:%s", CONFIG.ACTIVEMQ.HOST, CONFIG.ACTIVEMQ.PORT)

	return stomp.Dial("tcp", sAddr,
		stomp.ConnOpt.Login(CONFIG.ACTIVEMQ.USER, CONFIG.ACTIVEMQ.PASS),
		stomp.ConnOpt.HeartBeat(
			time.Duration(CONFIG.ACTIVEMQ.TIMEOUT)*time.Millisecond,
			time.Duration(CONFIG.ACTIVEMQ.TIMEOUT)*time.Millisecond,
		),
	)
}
