package cmd

import (
	"context"
	bootstrap "example/bootstrap"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oTcpCommand = &cobra.Command{
	Use:   "tcp",
	Short: "啟動 TCP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		// 收到中斷/終止訊號時 ctx 會被取消，Tcp.Serve 內部監聽 ctx.Done() 自己關掉 listener，
		// 不是靠 process 被系統強制殺掉才釋放 port。
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitTcpContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		oTcpRouter := register.TcpInit(oContainer)

		fmt.Println("TCP Server Start :" + bootstrap.CONFIG.SERVICES.TCP.PORT)

		if err := oTcpRouter.Serve(ctx, ":"+bootstrap.CONFIG.SERVICES.TCP.PORT); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oTcpCommand)
}
