package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
)

var oWebsocketCommand = &cobra.Command{
	Use:   "websocket",
	Short: "啟動 websocket 服務",
	Run: func(cmd *cobra.Command, args []string) {
		// 收到中斷/終止訊號時 ctx 會被取消，這個 ctx 會一路傳進 router，
		// 讓每條已經 upgrade 的連線自己決定要不要關掉；http.Server 另外靠 Shutdown
		// 停止 accept 新連線、關掉 listener，ListenAndServe() 才會正常返回並釋放 port，
		// 不是靠 process 被系統強制殺掉才釋放。
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitWebsocketContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		oRouter := register.WebsocketInit(oContainer)
		oRouter.Serve(ctx)

		oWebsocketServer := &http.Server{
			Addr: ":" + bootstrap.CONFIG.SERVICES.WEBSOCKET.PORT,
		}

		go func() {
			<-ctx.Done()
			oWebsocketServer.Shutdown(context.Background())
		}()

		if err := oWebsocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oWebsocketCommand)
}
