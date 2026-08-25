package register

import (
	"encoding/json"
	container "example/container"
	"fmt"
	"log"
	"net/http"

	pkg "example/pkg"

	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*

          open                                                                        完成
		  ping        / pong

	event:
	      connect    / conntected                                  reply ✅           完成
	      disconnect / disconnected 不需要                          reply ❌           完成
		  heartbeat  / heartbeated                                 reply ✅
		  subscribe  / subscribed                                  reply ✅

		  authenticate/ authenticated
	      presence, presence-stats, history,                       reply ✅           可取消
		  rpc        / rpc-ack                                     reply ✅

		  message    / messaged 不需要                              reply ❌
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅
		  unsubscribe / unsubscribed 不需要                        reply ❌

    method:
	value:


*/

func WebsocketInit(oContainer *container.WebsocketContainer) *http.ServeMux {

	oEventer := pkg.NewWebsocketEventer(websocket.Upgrader{
		CheckOrigin: func(oRequest *http.Request) bool {
			return true
		},
	})

	var (
		pointerToUuid    = hashmap.New[string, string]()
		uuidToConnection = hashmap.New[string, *websocket.Conn]()

		uuidToChannels = hashmap.New[string, string]()
		channelToConns = hashmap.New[string, string]()
	)

	_ = uuidToChannels
	_ = channelToConns

	oEventer.OnOpen(func(oConn *websocket.Conn) {
		log.Println("OnOpen:", oConn.RemoteAddr())

		sUuid := uuid.New().String()
		sPointer := fmt.Sprintf("%p", oConn)

		if _, bFound := uuidToConnection.Get(sUuid); bFound {
			log.Println("duplicate id, disconnect:", sUuid, oConn.RemoteAddr())
			oConn.Close()
			return
		}

		uuidToConnection.Set(sUuid, oConn)
		pointerToUuid.Set(sPointer, sUuid)

	})

	oEventer.OnConnect(func(oConn *websocket.Conn) {
		sPointer := fmt.Sprintf("%p", oConn)

		sUuId, _ := pointerToUuid.Get(sPointer)

		aByteJson, oErr := json.Marshal(struct {
			Type   string `json:"type"`
			Result string `json:"result"`
		}{
			Type:   "connect-ack",
			Result: sUuId,
		})

		if oErr != nil {
			log.Println("json marshal error:", oErr)
			return
		}

		oConn.WriteMessage(websocket.TextMessage, aByteJson)

	})

	oEventer.OnDisconnect(func(oConn *websocket.Conn) {
		// disconnect 是收不到 uuid 的
		sConnKey := fmt.Sprintf("%p", oConn)

		sId, _ := pointerToUuid.Get(sConnKey)
		pointerToUuid.Del(sConnKey)
		uuidToConnection.Del(sId)

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
