package bootstrap

import (
	"context"
	"fmt"

	"github.com/go-zeromq/zmq4"
)

// ZeroMQ 沒有像 amqp/nats 那樣統一的「連線」物件，每種 socket pattern
// 都要各自建立。這裡先開一個 PUB socket 讓其他服務可以 Dial 進來訂閱，
// 跟 amqp.go／nats.go 扮演的「對外廣播」角色一致。
func NewZeromq() (zmq4.Socket, error) {
	oSocket := zmq4.NewPub(context.Background())

	sEndpoint := fmt.Sprintf("tcp://%s:%s", CONFIG.ZEROMQ.HOST, CONFIG.ZEROMQ.PORT)

	if err := oSocket.Listen(sEndpoint); err != nil {
		return nil, err
	}

	return oSocket, nil
}
