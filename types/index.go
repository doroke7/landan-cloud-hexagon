package types

import "encoding/json"

type HttpRequestBody struct {
	V string `json:"v" form:"v" binding:"required"`
}

// TcpRequest 是 client 送給 server 的內容，method 用來給 Tcp 分發到對應的 handler。
// RequestId 由 client 產生，server 會原樣把它放進對應的 TcpResponse 裡回傳——
// 這樣一條連線上可以同時「掛」好幾個還沒回應的 request，client 端靠 RequestId
// 把亂序回來的 response 配對回正確的呼叫方，不需要每個並發請求各自佔一條連線。
type TcpRequest struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Method    string `json:"method"`
	Param     string `json:"param"`
}

// TcpResponse 是 server 回給 client 的內容，跟 pkg.Response 的 code/message/result 是同一套慣例。
type TcpResponse struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Result    any    `json:"result"`
}

// WebsocketRequest 是 client 送給 server 的內容，Method 用來給 Websocket 分發到對應的 handler，
// 跟 TcpRequest 是同一套慣例：RequestId 由 client 產生，server 原樣放進對應的
// WebsocketResponse 裡回傳（見 pkg.WebsocketRouter.serveConn），client 端靠它把回應
// 配對回正確的呼叫方，不用假設「先送先回」——這也是 sample/websocket/client.js 那個
// requestId -> callback Map 的做法。
//
// Param 用 json.RawMessage 而不是固定的 string，是因為每個 method 需要的參數形狀都不一樣
// （例如 SignIn 需要 {name, password}，可能有的 method 只要一個 id）；WebsocketRouter 只管
// 路由，不管參數長什麼樣，交給各自的 handler 自己 json.Unmarshal 成自己定義的參數 struct。
type WebsocketRequest struct {
	RequestId int    `json:"request-id"`
	Type      string `json:"type"` // event: 事件消息（client 發起的呼叫）
	C         string `json:"c"`
	P         string `json:"p"`
	S         string `json:"s"`
	O         string `json:"o"`

	Header json.RawMessage `json:"header"`
	Code   int             `json:"code"`
	Method string          `json:"method"`
	Param  json.RawMessage `json:"param"`
	Search json.RawMessage `json:"search"`
	Option json.RawMessage `json:"option"`
}

// WebsocketResponse 是 server 回給 client 的內容，跟 TcpResponse 的 code/message/result 是同一套慣例。
// Result 用 json.RawMessage 跟 Param 是同一個道理：每個 method 回傳的資料形狀都不一樣，
// WebsocketRouter 不管內容長什麼樣，由各自的 handler 自己 json.Marshal 成 RawMessage。
type WebsocketResponse struct {
	RequestId int    `json:"request-id"`
	Type      string `json:"type"` // ack：回傳消息，且需要客戶端 ack；normal：回傳消息客戶端不用 ack；none：不會回傳
	C         string `json:"c"`
	M         string `json:"m"`
	R         string `json:"r"`

	Header  json.RawMessage `json:"header"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// WebsocketConn 是 handler／middleware 拿到的連線操作介面，讓它們可以在處理這次
// request 的同時，直接對目前這條連線做事——目前只有 Push（主動推一筆 "event" 訊息
// 回這條連線，不算這次呼叫的 ack，client 端不會拿它去對應 RequestId 的 callback）。
// 實作在 pkg.WebsocketRouter 內部（不對外公開型別），以參數傳遞，不掛在
// WebsocketRequest 上。
type WebsocketConn interface {
	Push(sMethod string, oParam any) error
}

// WebsocketNextFunc 是呼叫鏈中「下一層」的簽名，最後一層對應到實際註冊的 WebsocketHandlerFunc。
type WebsocketNextFunc func(oConn WebsocketConn, oReq WebsocketRequest) WebsocketResponse

// WebsocketMiddlewareFunc 職責跟 gin.HandlerFunc 一樣：包住 fnNext，可以在業務邏輯前後插
// 自己的邏輯，呼叫 fnNext(oConn, oReq) 對應 gin.Context.Next()；不呼叫 fnNext、直接自己回
// 一個 WebsocketResponse 就等於 gin 的 Abort()。
type WebsocketMiddlewareFunc func(oConn WebsocketConn, oReq WebsocketRequest, fnNext WebsocketNextFunc) WebsocketResponse
