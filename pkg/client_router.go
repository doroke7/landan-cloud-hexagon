package pkg

import (
	"context"
	"log"
)

// ClientHandlerFunc 是一個 client-side 訂閱方法的簽名，職責跟 ConsumerHandlerFunc 一樣，
// 用來被 ClientRouter 依序註冊、並行啟動。
type ClientHandlerFunc func(ctx context.Context) error

// ClientRouter 職責跟 ConsumerRouter 一樣：只負責收集多個 client-side 訂閱方法，
// Serve 時每個各自在自己的 goroutine 執行，互不干擾。
type ClientRouter struct {
	handlers []ClientHandlerFunc
}

func NewClientRouter() *ClientRouter {
	return &ClientRouter{}
}

// Handle 註冊一個 client-side 訂閱方法，用法跟 ConsumerRouter.HandleFunc 一樣。
func (oSelf *ClientRouter) Handle(fnHandler ClientHandlerFunc) {
	oSelf.handlers = append(oSelf.handlers, fnHandler)
}

func (oSelf *ClientRouter) Serve(ctx context.Context) error {
	for _, fnHandler := range oSelf.handlers {
		go func(fn ClientHandlerFunc) {
			if err := fn(ctx); err != nil {
				log.Printf("client stopped: %v", err)
			}
		}(fnHandler)
	}

	<-ctx.Done()
	return nil
}
