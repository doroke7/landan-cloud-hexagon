package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var oUpgrader = websocket.Upgrader{
	CheckOrigin: func(oRequest *http.Request) bool {
		return true
	},
}

/*
	type: connect,                                                 reply ✅
	      disconnect,                                              reply ❌
		  heartbeat                                                reply ✅
		  subscribe,                                               reply ✅

	      presence, presence-stats, history,                       reply ✅
		  rpc,                                                     reply ✅

		  message,                                                 reply ❌
		  publish,                                                 reply ✅ + broadcast ✅

		  refresh,                                                 reply ✅
		  sub-refresh,                                             reply ✅
		  unsubscribe,                                             reply ❌

    method:
	value:


*/

func fnHandler(oWriter http.ResponseWriter, oRequest *http.Request) {

	oConn, oErr := oUpgrader.Upgrade(oWriter, oRequest, nil)

	if oErr != nil {
		log.Println(oErr)
		return
	}

	defer oConn.Close()

	// 讀取 client 訊息
	go func() {

		for {

			iMessageType, aMsg, oErr := oConn.ReadMessage()

			if oErr != nil {
				log.Println("read error:", oErr)
				return
			}

			// echo 回去
			oErr = oConn.WriteMessage(iMessageType, aMsg)

			if oErr != nil {
				return
			}

		}

	}()

	// server 主動推送
	oTicker := time.NewTicker(5 * time.Second)

	defer oTicker.Stop()

	for range oTicker.C {

		oErr := oConn.WriteMessage(websocket.TextMessage, []byte("server ping"))

		if oErr != nil {
			return
		}

	}

}

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 Websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {

		http.HandleFunc("/ws", fnHandler)

		log.Fatal(http.ListenAndServe(":8080", nil))

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
