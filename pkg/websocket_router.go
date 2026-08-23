package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type WebsocketHandlerFunc func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse

type websocketRouteNode struct {
	children    map[string]*websocketRouteNode
	middlewares []types.WebsocketMiddlewareFunc // 這個節點自己註冊的 middleware，不含祖先節點的
	handler     WebsocketHandlerFunc            // 這個節點自己註冊的 handler，沒有就是 nil
	registered  bool                            // 這個節點是不是真的被 Group() 註冊過，還是只是路過的中繼節點
}

func newWebsocketRouteNode() *websocketRouteNode {
	return &websocketRouteNode{children: make(map[string]*websocketRouteNode)}
}

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

type WebsocketGroup struct {
	router *WebsocketRouter
	node   *websocketRouteNode
	prefix string // 從 root 累加下來的完整路徑，組 HandleFunc 的 method key 用
}

func (oSelf *WebsocketGroup) Use(fnMiddlewares []types.WebsocketMiddlewareFunc) *WebsocketGroup {
	oSelf.node.middlewares = append(oSelf.node.middlewares, fnMiddlewares...)
	oSelf.node.registered = true
	return oSelf
}

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

	oCtx, fnCancel := context.WithCancel(oContext)
	defer fnCancel()

	go func() {
		<-oCtx.Done()
		oConn.Close()
	}()
	oConn.SetReadLimit(websocketMaxMessageSize)
	oConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	oConn.SetPongHandler(func(string) error {
		return oConn.SetReadDeadline(time.Now().Add(websocketPongWait))
	})

	var oWriteMu sync.Mutex
	go oSelf.ping(oCtx, oConn, &oWriteMu)

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
			oRes, bSeen := oSeenResponses[oReq.Id]
			oSeenMu.Unlock()

			fmt.Println("oReq=", oReq)

			if !bSeen {
				oRes = oSelf.dispatch(oWebsocketConn, oReq)
				oRes.Id = oReq.Id

				if oRes.Type == "" {
					oRes.Type = "normal"
				}

				oSeenMu.Lock()
				oSeenResponses[oReq.Id] = oRes
				oSeenMu.Unlock()
			}

			if oRes.Type == "none" {
				return
			}

			oWriteMu.Lock()
			defer oWriteMu.Unlock()

			if err := oConn.WriteJSON(oRes); err != nil {
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

	if oRequest.Type == "heartbeat" {
		return types.WebsocketResponse{Type: oRequest.Type, Code: 1, Message: ""}
	}

	if oRequest.Type == "event" {
		fnNext := oSelf.chain(oRequest.Method)
		if fnNext == nil {
			return types.WebsocketResponse{Code: -1, Message: errors.New("websocket: method not found").Error()}
		}
		return fnNext(oConn, oRequest)

	}

	if oRequest.Type == "broacast" {

		return types.WebsocketResponse{Type: oRequest.Type, Code: 1, Message: "尚未支持廣播"}

	}

	// Type: hearbeat ： 心跳
	// Type: event: 一般訊息，之後需要接受 server 來的 ack
	// Type: broacast： 廣播訊息 （client 傳給 server 後。 server loop 傳給其他 client， 並且透過 NATS）
	// Type: ack: 基本上只做 client -> server 丟消息後， server 過來的 ack

	return types.WebsocketResponse{Type: oRequest.Type, Code: 1, Message: "尚未支持"}
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
