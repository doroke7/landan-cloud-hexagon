package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	types "example/types"
)

const (
	websocketMaxMessageSize = 1 << 12 // 4KB，跟 TcpRouter 的 tcpMaxBodyLength 對稱，擋住異常/惡意的超大訊息把記憶體打爆
	websocketPongWait       = 60 * time.Second
	websocketPingPeriod     = websocketPongWait * 9 / 10 // 要比 pongWait 短，才能在逾時前送出下一個 ping
	websocketWriteWait      = 10 * time.Second
)

var ErrWebsocketMethodNotFound = errors.New("websocket: method not found")

type WebsocketHandlerFunc func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse

type WebsocketRouter struct {
	prefix       string
	upgrader     websocket.Upgrader
	routes       map[string]websocketRoute
	middlewares  []types.WebsocketMiddlewareFunc
	noMethodFunc types.WebsocketMiddlewareFunc
}

type websocketRoute struct {
	handler WebsocketHandlerFunc
	group   *WebsocketGroup
}

func NewWebsocketRouter(sPrefix string) *WebsocketRouter {
	return &WebsocketRouter{
		prefix: sPrefix,
		routes: make(map[string]websocketRoute),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (oSelf *WebsocketRouter) HandleFunc(sMethod string, fnHandler WebsocketHandlerFunc) *WebsocketRouter {
	oSelf.routes[sMethod] = websocketRoute{handler: fnHandler}
	return oSelf
}

func (oSelf *WebsocketRouter) Use(fnMiddlewares []types.WebsocketMiddlewareFunc) *WebsocketRouter {
	oSelf.middlewares = append(oSelf.middlewares, fnMiddlewares...)
	return oSelf
}

func (oSelf *WebsocketRouter) Group(sPrefix string) *WebsocketGroup {
	return &WebsocketGroup{router: oSelf, prefix: sPrefix}
}

// WebsocketGroup 用法跟 gin.RouterGroup 一樣：Use() 加這個 group 專屬的 middleware，
// HandleFunc() 用 group 的 prefix + method name 註冊到共用的 router 上。
type WebsocketGroup struct {
	router      *WebsocketRouter
	prefix      string
	middlewares []types.WebsocketMiddlewareFunc
}

// Use 註冊只套用到「這個 group 底下的 method」的 middleware，疊在 WebsocketRouter.Use
// 全局 middleware 之後、handler 之前。吃 slice 而不是 variadic，跟 WebsocketRouter.Use
// 保持同一套簽名。
func (oSelf *WebsocketGroup) Use(fnMiddlewares []types.WebsocketMiddlewareFunc) *WebsocketGroup {
	oSelf.middlewares = append(oSelf.middlewares, fnMiddlewares...)
	return oSelf
}

func (oSelf *WebsocketGroup) HandleFunc(sMethod string, fnHandler WebsocketHandlerFunc) *WebsocketGroup {
	oSelf.router.routes[oSelf.prefix+sMethod] = websocketRoute{handler: fnHandler, group: oSelf}
	return oSelf
}

// NoMethod 註冊「method 不存在」時要包住預設回應的 middleware，用法跟 gin.Engine.NoRoute
// 一樣：fnNext 是內建的 ErrWebsocketMethodNotFound 回應，NoMethod 可以在外面包一層自己的
// 邏輯（記錄、覆寫訊息……），不呼叫 fnNext 就等於自己決定要回什麼。沒註冊就直接用內建回應。
func (oSelf *WebsocketRouter) NoMethod(fnMiddleware types.WebsocketMiddlewareFunc) *WebsocketRouter {
	oSelf.noMethodFunc = fnMiddleware
	return oSelf
}

// Serve 把這個 router 註冊到全局 http（http.DefaultServeMux）的 prefix 路徑上；
// ctx 取消時（優雅關機），每條已經 upgrade 的連線都會被主動關掉，不用等 client 自己斷線。
func (oSelf *WebsocketRouter) Serve(oContext context.Context) {
	http.HandleFunc(oSelf.prefix, func(oResponseWriter http.ResponseWriter, oResquest *http.Request) {
		oSelf.serveConn(oContext, oResponseWriter, oResquest)
	})
}

func (oSelf *WebsocketRouter) serveConn(oContext context.Context, oResponseWriter http.ResponseWriter, oResquest *http.Request) {

	oConn, err := oSelf.upgrader.Upgrade(oResponseWriter, oResquest, nil)

	if err != nil {
		log.Printf("websocket: upgrade failed: %v", err)
		return
	}

	defer oConn.Close()

	// 用全局 ctx 衍生一個連線等級的子 ctx：全局 ctx 取消（優雅關機）或這條連線自己結束時
	// 都要能讓下面的 watcher/ping goroutine 退出，不然每條連線都會卡著永遠不返回的
	// goroutine，直到整個服務關機才釋放——是明確的 goroutine 洩漏。
	oCtx, fnCancel := context.WithCancel(oContext)
	defer fnCancel()

	go func() {
		<-oCtx.Done()
		oConn.Close()
	}()

	// SetReadLimit 擋住異常/惡意的超大訊息；pong handler 每收到一次 client 的 pong
	// 就把 read deadline 往後延，client 斷線或卡死超過 pongWait 沒回應，
	// 下面的 ReadJSON 就會因為逾時出錯、跳出迴圈，連線才不會無限期占著。
	oConn.SetReadLimit(websocketMaxMessageSize)
	oConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	oConn.SetPongHandler(func(string) error {
		return oConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	})

	// oWriteMu 保護「往同一個 oRawConn 寫東西」這個動作：gorilla websocket 規定同一條
	// 連線同時間只能有一個 goroutine 在寫，ping goroutine 的 WriteMessage(Ping)、下面回
	// ack、跟 handler 透過 oConn.Push 主動推播，三方都要序列化，跟 TcpRouter.serveConn
	// 的 oWriteMu 是同一個道理。
	var oWriteMu sync.Mutex
	go oSelf.ping(oCtx, oConn, &oWriteMu)

	// oConn 是傳給 handler／middleware 的 types.WebsocketConn 實作，整條連線只需要
	// 一份，不用每個 request 各自建一個。
	oWebsocketConn := NewWebsocketConn(oConn, &oWriteMu)

	var oSeenMu sync.Mutex
	oSeenResponses := make(map[int]types.WebsocketResponse)

	for {
		var oReq types.WebsocketRequest
		if err := oConn.ReadJSON(&oReq); err != nil {
			return
		}

		// 多路復用
		go func(oReq types.WebsocketRequest) {

			oSeenMu.Lock()
			oResp, bSeen := oSeenResponses[oReq.RequestId]
			oSeenMu.Unlock()

			if !bSeen {
				oResp = oSelf.dispatch(oWebsocketConn, oReq)
				oResp.RequestId = oReq.RequestId

				// Type 由 handler 自己決定（ack/normal/none），router 只根據這個欄位
				// 決定要不要真的送出去；handler 沒設就預設 "normal"（多數情況都不需要
				// client 額外回 ack，也不是完全不回應）。
				if oResp.Type == "" {
					oResp.Type = "normal"
				}

				oSeenMu.Lock()
				oSeenResponses[oReq.RequestId] = oResp
				oSeenMu.Unlock()
			}

			if oResp.Type == "none" {
				return
			}

			oWriteMu.Lock()
			defer oWriteMu.Unlock()

			if err := oConn.WriteJSON(oResp); err != nil {
				log.Printf("websocket: write failed: method=%s err=%v", oReq.Method, err)
			}
		}(oReq)
	}
}

func NewWebsocketConn(oConn *websocket.Conn, oWriteMu *sync.Mutex) *websocketConn {
	return &websocketConn{
		conn:    oConn,
		writeMu: oWriteMu,
	}
}

type websocketConn struct {
	conn    *websocket.Conn
	writeMu *sync.Mutex
}

func (oSelf *websocketConn) Push(sMethod string, oParam any) error {
	aParam, err := json.Marshal(oParam)
	if err != nil {
		return err
	}

	oSelf.writeMu.Lock()
	defer oSelf.writeMu.Unlock()

	return oSelf.conn.WriteJSON(types.WebsocketRequest{
		Type:   "event",
		Method: sMethod,
		Param:  aParam,
	})
}

func (oSelf *WebsocketRouter) dispatch(oConn types.WebsocketConn, oRequest types.WebsocketRequest) types.WebsocketResponse {

	oRoute, ok := oSelf.routes[oRequest.Method]

	if !ok {
		fnNotFound := types.WebsocketNextFunc(func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
			return types.WebsocketResponse{Code: -1, Message: ErrWebsocketMethodNotFound.Error()}
		})

		if oSelf.noMethodFunc == nil {
			return fnNotFound(oConn, oRequest)
		}
		return oSelf.noMethodFunc(oConn, oRequest, fnNotFound)
	}

	return oSelf.chain(oRoute)(oConn, oRequest)
}

func (oSelf *WebsocketRouter) chain(oRoute websocketRoute) types.WebsocketNextFunc {
	fnNext := types.WebsocketNextFunc(oRoute.handler)

	aMiddlewares := oSelf.middlewares
	if oRoute.group != nil {
		aMiddlewares = append(append([]types.WebsocketMiddlewareFunc{}, oSelf.middlewares...), oRoute.group.middlewares...)
	}

	for i := len(aMiddlewares) - 1; i >= 0; i-- {
		fnMiddleware := aMiddlewares[i]
		fnCurrentNext := fnNext

		fnNext = func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
			return fnMiddleware(oConn, oReq, fnCurrentNext)
		}
	}

	return fnNext
}

// ping 定期送 PingMessage 維持連線存活；client 沒有在 pongWait 內回應的話，
// 上面 ReadJSON 的 read deadline 會到期出錯，連線就會被 serveConn 清掉。
func (oSelf *WebsocketRouter) ping(ctx context.Context, oConn *websocket.Conn, oWriteMu *sync.Mutex) {
	oTicker := time.NewTicker(websocketPingPeriod)
	defer oTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-oTicker.C:
			oWriteMu.Lock()
			oConn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
			err := oConn.WriteMessage(websocket.PingMessage, nil)
			oWriteMu.Unlock()

			if err != nil {
				return
			}
		}
	}
}
