package cmd

import (
	"context"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"

	"example/pkg"

	"github.com/centrifugal/centrifuge"
	"github.com/spf13/cobra"
)

var oCentrifugeCommand = &cobra.Command{
	Use:   "centrifuge",
	Short: "啟動 centrifuge 服務",
	Run: func(cmd *cobra.Command, args []string) {

		oNats, oErr := nats.Connect(
			"nats://localhost:4222",
		)

		if oErr != nil {
			panic(oErr)
		}

		defer oNats.Close()

		oNats.Publish(
			"chat.room1",
			[]byte(
				`{"msg":"hello"}`,
			),
		)

		oBroker := pkg.NewNATSBroker(
			oNats,
		)

		oNode, oErr := centrifuge.New(
			centrifuge.Config{
				LogLevel: centrifuge.LogLevelDebug,
			},
		)

		if oErr != nil {
			panic(oErr)
		}

		oNode.SetBroker(oBroker)

		// Connecting：驗證/接受連線，這個 demo 不做任何驗證，一律接受成匿名連線。
		oNode.OnConnecting(
			func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
				return centrifuge.ConnectReply{}, nil
			},
		)

		// Connect：連線建立完成後才能拿到 *centrifuge.Client，在這裡掛 client 層級的事件、訂閱 channel。
		oNode.OnConnect(
			func(client *centrifuge.Client) {

				log.Println(
					"connected:",
					client.ID(),
				)

				// 斷線
				client.OnDisconnect(
					func(e centrifuge.DisconnectEvent) {

						log.Println(
							"disconnect:",
							client.ID(),
						)

					},
				)

				// client 訂閱 channel
				if err := client.Subscribe("chat:room1"); err != nil {
					log.Println("subscribe error:", err)
				}
			},
		)

		// 啟動 node
		if err := oNode.Run(); err != nil {
			panic(err)
		}

		fnHandler := centrifuge.NewWebsocketHandler(
			oNode,
			centrifuge.WebsocketConfig{},
		)

		http.Handle(
			"/connection/websocket",
			fnHandler,
		)

		go func() {

			for {

				// 模擬 broadcast

				oNode.Publish(
					"chat:room1",
					[]byte(
						`{"message":"hello"}`,
					),
				)

			}

		}()

		log.Println(
			"listen :8000",
		)

		http.ListenAndServe(
			":8000",
			nil,
		)
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCentrifugeCommand)
}
