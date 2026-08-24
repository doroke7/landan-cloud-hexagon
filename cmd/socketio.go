package cmd

import (
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

/*
	event: connect,                                                reply ✅
	       disconnect,                                              reply ❌
		   message,                                                 reply ✅ (echo)
		   broadcast,                                                reply ✅ + broadcast ✅

    server -> client:
	       server ping (每 5 秒推播一次)                              reply ✅
*/

func fnSocketIOHandler() *socketio.Server {

	oServer := socketio.NewServer(nil)

	oServer.OnConnect("/", func(oConn socketio.Conn) error {

		oConn.SetContext("")
		log.Println("connected:", oConn.ID())

		return nil

	})

	oServer.OnEvent("/", "message", func(oConn socketio.Conn, sMsg string) string {

		// echo 回去
		return sMsg

	})

	oServer.OnEvent("/", "broadcast", func(oConn socketio.Conn, sMsg string) {

		oServer.BroadcastToNamespace("/", "broadcast", sMsg)

	})

	oServer.OnError("/", func(oConn socketio.Conn, oErr error) {
		log.Println("error:", oErr)
	})

	oServer.OnDisconnect("/", func(oConn socketio.Conn, sReason string) {
		log.Println("disconnected:", oConn.ID(), sReason)
	})

	return oServer

}

func fnSocketIOMain() {

	oServer := fnSocketIOHandler()

	go oServer.Serve()
	defer oServer.Close()

	// server 主動推送
	go func() {

		oTicker := time.NewTicker(5 * time.Second)
		defer oTicker.Stop()

		for range oTicker.C {
			oServer.BroadcastToNamespace("/", "server ping", "server ping")
		}

	}()

	http.Handle("/socket.io/", oServer)

	http.ListenAndServe(":8081", nil)

}
