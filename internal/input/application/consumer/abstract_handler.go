package consumer

import (
	amqp "github.com/rabbitmq/amqp091-go"

	helper "example/internal/helper"
)

// AbstractHandler 放 consumer 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / http）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*helper.AesHelper
	Conn *amqp.Connection
}

func NewAbstractHandler(oAesHelper *helper.AesHelper, oConn *amqp.Connection) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
		Conn:      oConn,
	}
}
