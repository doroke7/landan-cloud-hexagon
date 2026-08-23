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

/*
WebsocketRouter 的 group 模仿 GrpcRouter（見 pkg/grpc_router.go）：prefix 用「.」「/」
分段建成一棵樹，dispatch 時沿著 oRequest.Method 的每一段往下走，沿路每經過一個真的被
Group() 註冊過的節點，就把該節點自己的 middleware 疊加上去——子 group 會自動繼承所有
祖先 group 的 middleware，再疊加自己的，順序是「祖先在前、自己在後」，且跟 Group()
呼叫的先後順序完全無關。tokenize() 直接複用 grpc_router.go 那個（同一個 package）。

沒有獨立的「全局 middlewares」概念：root 本身就是樹的一個節點，想要「不管打哪個 method
都要跑」的全局 middleware，直接 Group("", 全局middleware...) 註冊在 root 上即可，
dispatch 時一定會先經過 root，效果一樣，只是統一成同一套機制，不用額外的欄位。

handler 也直接掛在對應的樹節點上（不再另外開一個 map[string]handler）：HandleFunc
跟 Group 用同一套 tokenize 建樹，method 裡的「.」「/」視為同一種分隔符——這樣
Group("Admin").HandleFunc("/Authentication/Authenticator.SignIn", ...) 跟 client 端
送出的 "Admin.Authentication.Authenticator.SignIn"（分隔符不同）才能對到同一個節點；
如果 handler 另外存一份用原始字串當 key 的 map，兩種分隔符寫法會被當成不同的 key，
永遠查不到。
*/
type websocketRouteNode struct {
	children    map[string]*websocketRouteNode
	middlewares []types.WebsocketMiddlewareFunc // 這個節點自己註冊的 middleware，不含祖先節點的
	handler     WebsocketHandlerFunc            // 這個節點自己註冊的 handler，沒有就是 nil
	registered  bool                            // 這個節點是不是真的被 Group() 註冊過，還是只是路過的中繼節點
}

func newWebsocketRouteNode() *websocketRouteNode {
	return &websocketRouteNode{children: make(map[string]*websocketRouteNode)}
}

// websocketNodeFor 沿著 sPath 的每一段走過樹，不存在的節點沿路建起來，回傳最終的
// 葉節點；Group()、WebsocketGroup.Group()、HandleFunc() 都靠它定位/建樹，只是
// 起點不同（router 從 root 開始，group 從自己的節點開始）。
func websocketNodeFor(oNode *websocketRouteNode, sPath string) *websocketRouteNode {
	for _, sSegment := range tokenize(sPath) {
		oChild, bOk := oNode.children[sSegment]
		if !bOk {
			oChild = newWebsocketRouteNode()
			oNode.children[sSegment] = oChild
		}
		oNode = oChild
	}
	return oNode
}

type WebsocketRouter struct {
	prefix   string
	upgrader websocket.Upgrader
	root     *websocketRouteNode
}

func NewWebsocketRouter(sPrefix string) *WebsocketRouter {
	return &WebsocketRouter{
		prefix: sPrefix,
		root:   newWebsocketRouteNode(),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (oSelf *WebsocketRouter) HandleFunc(sMethod string, fnHandler WebsocketHandlerFunc) *WebsocketRouter {
	websocketNodeFor(oSelf.root, sMethod).handler = fnHandler
	return oSelf
}

// Group 每次 aaa.bbb.ccc 的 path 就生成 tree 的結構，並且在最後的節點掛上 middleware，
// 用法跟 pkg.GrpcRouter.Group 是同一套慣例。
func (oSelf *WebsocketRouter) Group(sPrefix string, aMiddlewares ...types.WebsocketMiddlewareFunc) *WebsocketGroup {
	oNode := websocketNodeFor(oSelf.root, sPrefix)
	oNode.middlewares = append(oNode.middlewares, aMiddlewares...)
	oNode.registered = true

	return &WebsocketGroup{router: oSelf, node: oNode, prefix: sPrefix}
}

// WebsocketGroup 用法跟 gin.RouterGroup 一樣：Use() 加這個 group 專屬的 middleware，
// HandleFunc() 用 group 的 prefix + method name 註冊到共用的 router 上；Group() 可以
// 在這個 group 底下再開子 group（children），一路巢狀下去，跟 GrpcRouter 的樹狀結構
// 是同一顆樹，只是進入點不同（從這個節點往下建，不是從 root）。
type WebsocketGroup struct {
	router *WebsocketRouter
	node   *websocketRouteNode
	prefix string // 從 root 累加下來的完整路徑，組 HandleFunc 的 method key 用
}

// Use 註冊只套用到「這個 group（含底下所有 children）」的 middleware，疊在祖先 group
// 的 middleware 之後、handler 之前。吃 slice 而不是 variadic，跟 WebsocketGroup.Group
// 的 variadic 簽名不同，是刻意維持跟舊版一樣的呼叫方式。
func (oSelf *WebsocketGroup) Use(fnMiddlewares []types.WebsocketMiddlewareFunc) *WebsocketGroup {
	oSelf.node.middlewares = append(oSelf.node.middlewares, fnMiddlewares...)
	oSelf.node.registered = true
	return oSelf
}

// Group 在這個 group 底下開一個子 group（children），繼承這個 group 的完整路徑，
// 子 group 的 middleware 只套用在它自己跟它底下的 method，不會影響同層的其他 group。
func (oSelf *WebsocketGroup) Group(sPrefix string, aMiddlewares ...types.WebsocketMiddlewareFunc) *WebsocketGroup {
	oNode := websocketNodeFor(oSelf.node, sPrefix)
	oNode.middlewares = append(oNode.middlewares, aMiddlewares...)
	oNode.registered = true

	return &WebsocketGroup{router: oSelf.router, node: oNode, prefix: oSelf.prefix + sPrefix}
}

func (oSelf *WebsocketGroup) HandleFunc(sMethod string, fnHandler WebsocketHandlerFunc) *WebsocketGroup {
	websocketNodeFor(oSelf.node, sMethod).handler = fnHandler
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
			oResp, bSeen := oSeenResponses[oReq.Id]
			oSeenMu.Unlock()

			if !bSeen {
				oResp = oSelf.dispatch(oWebsocketConn, oReq)
				oResp.Id = oReq.Id

				// Type 由 handler 自己決定（ack/normal/none），router 只根據這個欄位
				// 決定要不要真的送出去；handler 沒設就預設 "normal"（多數情況都不需要
				// client 額外回 ack，也不是完全不回應）。
				if oResp.Type == "" {
					oResp.Type = "normal"
				}

				oSeenMu.Lock()
				oSeenResponses[oReq.Id] = oResp
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
		Value:  aParam,
	})
}

func (oSelf *WebsocketRouter) dispatch(oConn types.WebsocketConn, oRequest types.WebsocketRequest) types.WebsocketResponse {
	fnNext := oSelf.chain(oRequest.Method)
	if fnNext == nil {
		return types.WebsocketResponse{Code: -1, Message: ErrWebsocketMethodNotFound.Error()}
	}

	return fnNext(oConn, oRequest)
}

// chain 沿著 sMethod 的每一段走過樹（跟 GrpcRouter.Build 內部查找同一套邏輯），一路
// 把「真的被 Group() 註冊過」的節點的 middleware 依序疊上去（祖先在前、自己在後），
// 走到底那個節點的 handler 就是最終要呼叫的方法；沒有 handler（純中繼節點，或
// method 根本沒註冊過）回傳 nil，交給呼叫端決定要怎麼回應「method not found」。
func (oSelf *WebsocketRouter) chain(sMethod string) types.WebsocketNextFunc {
	var aMiddlewares []types.WebsocketMiddlewareFunc

	oNode := oSelf.root
	if oNode.registered {
		aMiddlewares = append(aMiddlewares, oNode.middlewares...)
	}

	for _, sSegment := range tokenize(sMethod) {
		oChild, bOk := oNode.children[sSegment]
		if !bOk {
			return nil
		}
		oNode = oChild
		if oNode.registered {
			aMiddlewares = append(aMiddlewares, oNode.middlewares...)
		}
	}

	if oNode.handler == nil {
		return nil
	}

	fnNext := types.WebsocketNextFunc(oNode.handler)

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
