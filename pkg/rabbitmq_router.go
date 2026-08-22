package pkg

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumerHandlerFunc 是一個 queue 對應的處理方法，簽名統一，方便用 queue name 當 key 做路由。
type RabbitmqHandlerFunc func(msg amqp.Delivery)

// RabbitmqRouter 職責跟 http.ServeMux／grpc 的 service registry 一樣，
// 只負責把 queue name 對應到一個處理方法，不管 unmarshal／business 邏輯。
type RabbitmqRouter struct {
	conn   *amqp.Connection
	routes map[string]RabbitmqHandlerFunc
}

func NewRabbitmqRouter(oConn *amqp.Connection) *RabbitmqRouter {
	return &RabbitmqRouter{
		conn:   oConn,
		routes: make(map[string]RabbitmqHandlerFunc),
	}
}

// HandleFunc 註冊一個 queue name 對應的處理方法，用法跟 http.HandleFunc 一樣。
func (oSelf *RabbitmqRouter) HandleFunc(sQueue string, fnHandler RabbitmqHandlerFunc) {
	oSelf.routes[sQueue] = fnHandler
}

func (oSelf *RabbitmqRouter) Serve(ctx context.Context) error {
	oChannel, err := oSelf.conn.Channel()
	if err != nil {
		return err
	}
	defer oChannel.Close()

	for sQueue, fnHandler := range oSelf.routes {
		q, err := oChannel.QueueDeclare(sQueue, true, false, false, false, nil)
		if err != nil {
			return err
		}

		msgs, err := oChannel.Consume(q.Name, "", false, false, false, false, nil)
		if err != nil {
			return err
		}

		go consume(ctx, msgs, fnHandler)
	}

	<-ctx.Done()
	return nil
}

func consume(ctx context.Context, msgs <-chan amqp.Delivery, fnHandler RabbitmqHandlerFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			fnHandler(msg)
		}
	}
}
