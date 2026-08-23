package cmd

import (
	"context"
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

				// client 訂閱 channel
				if oErr := oClient.Subscribe("all"); oErr != nil {
					log.Println("subscribe error:", oErr)
				}
			},
		)

		// 啟動 node
		if oErr := oNode.Run(); oErr != nil {
			panic(oErr)
		}

		fnHandler := centrifuge.NewWebsocketHandler(
			oNode,
			centrifuge.WebsocketConfig{
				// CORS：centrifuge 預設用 sameHostOriginCheck，要求 Origin host 跟 request Host
				// 一致，跨源的瀏覽器前端會直接被拒絕，這裡跟 pkg/websocket_router.go 同一套慣例全部放行。
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
			oTicker := time.NewTicker(time.Second * 10)
			defer oTicker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-oTicker.C:
					if _, oErr := oNode.Publish("all", []byte(`{"message": "Hi!!"}`)); oErr != nil {
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
