package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
		oNode.OnConnecting(
			func(oCtx context.Context, oEvent centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
				return centrifuge.ConnectReply{}, nil
			},
		)

		// Connect：連線建立完成後才能拿到 *centrifuge.Client，在這裡掛 client 層級的事件、訂閱 channel。
		oNode.OnConnect(
			func(oClient *centrifuge.Client) {

				log.Println(
					"connected:",
					oClient.ID(),
				)

				// 斷線
				oClient.OnDisconnect(
					func(e centrifuge.DisconnectEvent) {

						log.Println(
							"disconnect:",
							oClient.ID(),
						)

					},
				)

				// client 訂閱 channel
				if oErr := oClient.Subscribe("chat:room1"); oErr != nil {
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
			centrifuge.WebsocketConfig{},
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

			for {

				select {
				case <-ctx.Done():
					return
				default:
				}

				// 模擬 broadcast
				oNode.Publish(
					"chat:room1",
					[]byte(
						`{"message":"hello"}`,
					),
				)

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
