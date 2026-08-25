package register

import (
	container "example/container"
	"fmt"
	"log"
	"net/http"

	pkg "example/pkg"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebsocketInit 只負責注入業務邏輯（連線註冊、log、echo），連線生命週期機制
// （upgrade、read loop、斷線偵測）交給 pkg.WebsocketEventer，跟 SocketioInit／
// CentrifugeInit 是同一套「通訊邏輯跟業務邏輯分開」的慣例。
func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oEventer := pkg.NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	// 連線註冊表：id 一律由 server 端在 OnConnect 生成，不能信任 client 自己宣稱的 id。
	// 用 cornelk/hashmap 而不是 map + Mutex／sync.Map——它是真正 lock-free 的實作，
	// 讀寫都不需要拿鎖。hashmap.Map 的 Key 限制只能是數字／字串（不能是指標），
	// oIDByConn 用 fmt.Sprintf("%p", oConn) 把指標位址轉成字串當 key 來繞過這個限制；
	// OnDisconnect 只拿得到 *websocket.Conn，要靠它反查回 id 才知道斷的是哪一個。
	var (
		connectionMap   = hashmap.New[string, *websocket.Conn]()
		oIdToConnection = hashmap.New[string, string]()
	)

	oEventer.OnConnect(func(oConn *websocket.Conn) {
		sId := uuid.New().String()
		sConnKey := fmt.Sprintf("%p", oConn)

		connectionMap.Set(sId, oConn)
		oIdToConnection.Set(sConnKey, sId)

		log.Println("connected:", sId, oConn.RemoteAddr())
	})
	oEventer.OnDisconnect(func(oConn *websocket.Conn) {
		sConnKey := fmt.Sprintf("%p", oConn)

		sId, _ := oIdToConnection.Get(sConnKey)
		oIdToConnection.Del(sConnKey)
		connectionMap.Del(sId)

		log.Println("disconnected:", sId, oConn.RemoteAddr())
	})
	oEventer.OnMessage(func(oConn *websocket.Conn, iMessageType int, aMsg []byte) {
		// echo 回去
		oConn.WriteMessage(iMessageType, aMsg)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/ws", oEventer)

	return oMux
}
