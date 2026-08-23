package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	pkg "example/pkg"
)

/*
   Channel:
   Type: heartbeat / heartbeat-ack
   Method:
   Value:



   Channel:
   Type: room / room-ack
   Method: join/leave
   Value: room-01


   Channel:
   Type: event / event-ack
   Method: Admin/Authentication/Authenticator/SignIn
   Value: {}


   Channel: room-01
   Type: message / message-ack
   Method: broast
   Value: {}

   Channel: room-01
   Type: message / messages-ack
   Method: send
   Value: {}



*/

var oCentrifugeCommand = &cobra.Command{
	Use:   "centrifuge",
	Short: "啟動 centrifuge 服務",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, oErr := container.InitCentrifugeContainer(ctx)
		if oErr != nil {
			log.Fatal(oErr)
		}
		defer oContainer.Nats.Close()

		oNode, oErr := centrifuge.New(
			centrifuge.Config{
				LogLevel: centrifuge.LogLevelDebug,
			},
		)

		if oErr != nil {
			panic(oErr)
		}

		oNode.SetBroker(oContainer.NatsBroker)

		// Connecting：驗證/接受連線，這個 demo 不做任何驗證，一律接受成匿名連線。
		// Credentials 一定要給值（UserID 留空即代表匿名），不然 centrifuge 會在
		// connectCmd 判斷 credentials == nil 直接以 bad request 斷線。
		oNode.OnConnecting(
			func(oCtx context.Context, oEvent centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
				return centrifuge.ConnectReply{
					Credentials: &centrifuge.Credentials{UserID: ""},
				}, nil
			},
		)

		// Connect：連線建立完成後才能拿到 *centrifuge.Client，在這裡掛 client 層級的事件、訂閱 channel。
		oNode.OnConnect(
			func(oClient *centrifuge.Client) {

				// 斷線
				oClient.OnDisconnect(func(oEvent centrifuge.DisconnectEvent) {})

				if oErr := oClient.Subscribe("all"); oErr != nil {
					log.Println("subscribe error:", oErr)
				}
				// if oErr := oClient.Subscribe("room-01"); oErr != nil {
				// 	log.Println("subscribe error:", oErr)
				// }

				/*
				   1. 所有的 client 過來的消息 ， OnPublish 都會handler
				*/
				oClient.OnPublish(func(oEvent centrifuge.PublishEvent, fnCallback centrifuge.PublishCallback) {
					log.Printf(
						"收到訊息 channel=%s method=%s value=%s\n",
						oEvent.Channel,
					)
					var oPayload struct {
						Method string          `json:"method"`
						Value  json.RawMessage `json:"value"`
					}
					if oErr := json.Unmarshal(oEvent.Data, &oPayload); oErr != nil {
						log.Println("publish payload 解析失敗:", oErr)
					} else {
						log.Printf(
							"收到訊息 channel=%s method=%s value=%s\n",
							oEvent.Channel,
							oPayload.Method,
							oPayload.Value,
						)
					}

					fnCallback(centrifuge.PublishReply{}, nil)
				})
			},
		)

		// 啟動 node
		if oErr := oNode.Run(); oErr != nil {
			panic(oErr)
		}

		fnHandler := centrifuge.NewWebsocketHandler(
			oNode,
			centrifuge.WebsocketConfig{

				CheckOrigin: func(oRequest *http.Request) bool { return true },
			},
		)

		oMux := http.NewServeMux()
		oMux.Handle("/connection/websocket", fnHandler)

		oCentrifugeServer := &http.Server{
			Addr:    ":" + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT,
			Handler: oMux,
		}

		go func() {
			<-ctx.Done()
			oNode.Shutdown(context.Background())
			oCentrifugeServer.Shutdown(context.Background())
		}()

		go func() {
			oTicker10 := time.NewTicker(time.Second * 10)
			oTicker5 := time.NewTicker(time.Second * 5)

			defer oTicker10.Stop()
			defer oTicker5.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-oTicker10.C:
					if _, oErr := oNode.Publish("all", []byte(`{"message": "Hi!!"}`)); oErr != nil {
						log.Println("publish error:", oErr)
					}
				case <-oTicker5.C:
					if _, oErr := oNode.Publish("room-01", []byte(`{"message": "你好"}`)); oErr != nil {
						log.Println("publish error:", oErr)
					}
				}
			}
		}()

		pkg.Logger(pkg.Default).Info("啟動 CENTRIFUGE 服務。 port: " + bootstrap.CONFIG.SERVICES.CENTRIFUGE.PORT)

		if oErr := oCentrifugeServer.ListenAndServe(); oErr != nil && oErr != http.ErrServerClosed {
			log.Fatal(oErr)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCentrifugeCommand)
}
