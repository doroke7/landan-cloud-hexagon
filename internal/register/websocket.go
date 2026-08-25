package register

import (
	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pkg "example/pkg"

	container "example/container"
	types "example/types"
)

/*

          open                                                                        完成
		  ping        / pong

	event:
	      connect    / conntected                                  reply ✅           完成
	      disconnect / disconnected 不需要                          reply ❌           完成
		  heartbeat  / heartbeated                                 reply ✅

		  authenticate/ authenticated                              reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------


	      presence, presence-stats, history,                       reply ✅           可取消
		  rpc        / rpced                                       reply ✅

		  subscribe  / subscribed                                  reply ✅

		  message    / messaged 不需要                              reply ❌
		  broadcast  / broadcasted                                 reply ✅ + broadcast ✅

		  refresh     /refreshed                                  reply ✅
		  sub-refresh / sub-refreshed                             reply ✅

--------------------------------------需要檢查是否 authenticated -----------------------------------------------

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
		pointerToUuid           = hashmap.New[string, string]()
		uuidToConnection        = hashmap.New[string, *websocket.Conn]()
		pointerToAuthentication = hashmap.New[string, bool]()

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

	oEventer.OnConnect(func(oConn *websocket.Conn, oWsReq *types.WebsocketRequest) {
		sPointer := fmt.Sprintf("%p", oConn)

		sUuId, _ := pointerToUuid.Get(sPointer)

		aByteMessage, oErr := json.Marshal(struct {
			Event  string `json:"event"`
			Result string `json:"result"`
			RId    string `json:"r_id"`
		}{
			Event:  "connected",
			Result: sUuId,
			RId:    oWsReq.RId,
		})

		if oErr != nil {
			log.Println("json marshal error:", oErr)
			return
		}

		oConn.WriteMessage(websocket.TextMessage, aByteMessage)

	})

	oEventer.OnAuthenticate(func(oConn *websocket.Conn, oWsReq *types.WebsocketRequest) bool {
		var oValue struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}

		if oErr := json.Unmarshal(oWsReq.Value, &oValue); oErr != nil {
			log.Println("json unmarshal error:", oErr)
		}

		bOk := oValue.Name == "admin" && oValue.Password == "123456"

		nCode := -1
		if bOk {
			nCode = 1
		}

		aByteMessage, oErr := json.Marshal(struct {
			Event string `json:"event"`
			Code  int    `json:"code"`
			RId   string `json:"r_id"`
		}{
			Event: "authenticated",
			Code:  nCode,
			RId:   oWsReq.RId,
		})

		if oErr != nil {
			log.Println("json marshal error:", oErr)
			return false
		}

		oConn.WriteMessage(websocket.TextMessage, aByteMessage)

		if bOk {
			sPointer := fmt.Sprintf("%p", oConn)
			pointerToAuthentication.Set(sPointer, true)
		}

		return bOk
	})

	// disconnect 是收不到 uuid 的

	oEventer.OnDisconnect(func(oConn *websocket.Conn) {
		sConnKey := fmt.Sprintf("%p", oConn)

		sId, _ := pointerToUuid.Get(sConnKey)
		pointerToUuid.Del(sConnKey)
		uuidToConnection.Del(sId)
		pointerToAuthentication.Del(sConnKey)

		log.Println("disconnected:", sId, oConn.RemoteAddr())
	})
	oEventer.OnMessage(func(oConn *websocket.Conn, iMessageType int, aMsg []byte) {
		sPointer := fmt.Sprintf("%p", oConn)

		if bAuthenticated, _ := pointerToAuthentication.Get(sPointer); !bAuthenticated {
			log.Println("not authenticated, ignore message:", oConn.RemoteAddr())
			return
		}

		// echo 回去
		oConn.WriteMessage(iMessageType, aMsg)
	})

	oMux := http.NewServeMux()
	oMux.Handle("/ws", oEventer)

	return oMux
}
