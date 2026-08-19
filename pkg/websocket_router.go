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

// WebsocketHandlerFunc 是一個 method 對應的處理方法，簽名統一，方便用 method name 當 key 做路由，
// 跟 pkg.TcpHandlerFunc 是同一套慣例。oConn 讓 handler 可以在處理這次 request 的同時，
// 直接對目前這條連線做事（例如主動 Push 一筆額外的 event），是獨立的參數，不掛在
// oReq 上面。
type WebsocketHandlerFunc func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse

// WebsocketRouter 職責跟 TcpRouter 一樣：只負責把 method name 對應到一個處理方法，不管
// unmarshal／business 邏輯；一條 websocket 連線 upgrade 完之後可以依序呼叫多個不同 method
// 的業務邏輯，不用每個業務各自佔一個 http path、各自重複處理 upgrade／心跳／優雅關機。
type WebsocketRouter struct {
	prefix       string
	upgrader     websocket.Upgrader
	routes       map[string]websocketRoute
	middlewares  []types.WebsocketMiddlewareFunc
	noMethodFunc types.WebsocketMiddlewareFunc
}

// websocketRoute 記著這個 method 註冊在哪個 group（沒有就是 nil，代表直接註冊在 router
// 上）；group 存指標而不是複製一份 middleware slice，是為了讓 dispatch 時能讀到
// group.Use() 之後才追加的 middleware，不用要求 Use() 一定要在 HandleFunc 之前呼叫。
type websocketRoute struct {
	handler WebsocketHandlerFunc
	group   *WebsocketGroup
}

func NewWebsocketRouter(sPrefix string) *WebsocketRouter {
	return &WebsocketRouter{
		prefix: sPrefix,
		routes: make(map[string]websocketRoute),
	}
}

// HandleFunc 註冊一個 method 對應的處理方法，用法跟 TcpRouter.HandleFunc 一樣。
func (oSelf *WebsocketRouter) HandleFunc(sMethod string, fnHandler WebsocketHandlerFunc) *WebsocketRouter {
	oSelf.routes[sMethod] = websocketRoute{handler: fnHandler}
	return oSelf
}

// Use 註冊套用到「這個 router 上所有 method」的 middleware，用法跟 gin.RouterGroup.Use
// 一樣：依註冊順序由外往內包住實際的 handler，第一個註冊的最外層，先執行、最後收尾。
// 吃 slice 而不是 variadic，跟 WebsocketGroup.Use 保持同一套簽名。
func (oSelf *WebsocketRouter) Use(fnMiddlewares []types.WebsocketMiddlewareFunc) *WebsocketRouter {
	oSelf.middlewares = append(oSelf.middlewares, fnMiddlewares...)
	return oSelf
}

// Group 對應 gin.RouterGroup：回傳一個共用同一個 WebsocketRouter（同一條連線、同一個
// dispatch table）的路由分群，method 名稱會自動加上 sPrefix；Group.Use() 疊加的
// middleware 只套用在這個 group 底下註冊的方法，跑在 router 全局 middleware 之後、
// handler 之前——一個 router 可以開多個 group（例如 admin/app/third 各自一個），
// 共用同一條連線跟心跳／優雅關機邏輯，各自的業務 middleware 互不影響。
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
func (oSelf *WebsocketRouter) Serve(ctx context.Context) {
	http.HandleFunc(oSelf.prefix, func(w http.ResponseWriter, r *http.Request) {
		oSelf.serveConn(ctx, w, r)
	})
}

func (oSelf *WebsocketRouter) serveConn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	oRawConn, err := oSelf.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket: upgrade failed: %v", err)
		return
	}
	defer oRawConn.Close()

	// 用全局 ctx 衍生一個連線等級的子 ctx：全局 ctx 取消（優雅關機）或這條連線自己結束時
	// 都要能讓下面的 watcher/ping goroutine 退出，不然每條連線都會卡著永遠不返回的
	// goroutine，直到整個服務關機才釋放——是明確的 goroutine 洩漏。
	oCtx, fnCancel := context.WithCancel(ctx)
	defer fnCancel()

	go func() {
		<-oCtx.Done()
		oRawConn.Close()
	}()

	// SetReadLimit 擋住異常/惡意的超大訊息；pong handler 每收到一次 client 的 pong
	// 就把 read deadline 往後延，client 斷線或卡死超過 pongWait 沒回應，
	// 下面的 ReadJSON 就會因為逾時出錯、跳出迴圈，連線才不會無限期占著。
	oRawConn.SetReadLimit(websocketMaxMessageSize)
	oRawConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	oRawConn.SetPongHandler(func(string) error {
		return oRawConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	})

	// oWriteMu 保護「往同一個 oRawConn 寫東西」這個動作：gorilla websocket 規定同一條
	// 連線同時間只能有一個 goroutine 在寫，ping goroutine 的 WriteMessage(Ping)、下面回
	// ack、跟 handler 透過 oConn.Push 主動推播，三方都要序列化，跟 TcpRouter.serveConn
	// 的 oWriteMu 是同一個道理。
	var oWriteMu sync.Mutex
	go oSelf.ping(oCtx, oRawConn, &oWriteMu)

	// oConn 是傳給 handler／middleware 的 types.WebsocketConn 實作，整條連線只需要
	// 一份，不用每個 request 各自建一個。
	oWebsocketConn := NewWebsocketConn(oRawConn, &oWriteMu)

	// oSeenResponses 快取「這個 RequestId 算過的回應」：script/websocket/client.js 的
	// emit() 如果一直沒收到回應會重送同一個 RequestId，如果每次重送都重新 dispatch，
	// 像 SignIn 這種有副作用（發 JWT、Push 事件）的 handler 就會被多執行一次——
	// 命中快取就直接回舊的 oResp，不重新跑 handler。只在這條連線活著的時候存在，
	// 連線關閉就跟著釋放，不會無限增長。
	var oSeenMu sync.Mutex
	oSeenResponses := make(map[int]types.WebsocketResponse)

	// 跟 TcpRouter.serveConn 一樣：read 迴圈本身是循序的（一次只解一個 request），
	// 但每個 request 都丟到自己的 goroutine 去 dispatch——這樣一個慢 method 不會卡住
	// 後面訊息的讀取，同一條連線可以多路復用，回應順序也可能跟收到的順序不一樣，
	// 這就是為什麼 RequestId 存在的意義（client 端靠它配對回正確的呼叫方，見
	// sample/websocket/client.js 的 requestId -> callback 對應）。
	for {
		var oReq types.WebsocketRequest
		if err := oRawConn.ReadJSON(&oReq); err != nil {
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

			if err := oRawConn.WriteJSON(oResp); err != nil {
				log.Printf("websocket: write failed: method=%s err=%v", oReq.Method, err)
			}
		}(oReq)
	}
}

func NewWebsocketConn(conn *websocket.Conn, writeMu *sync.Mutex) *websocketConn {
	return &websocketConn{
		conn:    conn,
		writeMu: writeMu,
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

func (oSelf *WebsocketRouter) dispatch(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
	oRoute, ok := oSelf.routes[oReq.Method]
	if !ok {
		fnNotFound := types.WebsocketNextFunc(func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
			return types.WebsocketResponse{Code: -1, Message: ErrWebsocketMethodNotFound.Error()}
		})

		if oSelf.noMethodFunc == nil {
			return fnNotFound(oConn, oReq)
		}
		return oSelf.noMethodFunc(oConn, oReq, fnNotFound)
	}

	return oSelf.chain(oRoute)(oConn, oReq)
}

// chain 把套用到這個 method 的 middleware 由外往內包住 handler，組成跟 gin
// Use()/Next() 一樣的洋蔥式呼叫鏈：router 全局 middleware 最外層，這個 method 所屬
// group 自己的 middleware 接著疊上去，最後才輪到 handler；router.middlewares 跟
// group.middlewares 都是每次 dispatch 時現讀，Use() 呼叫的先後順序不影響結果。
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
