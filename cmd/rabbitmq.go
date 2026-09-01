package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	registerRabbitmq "example/internal/register/rabbitmq"
	pkgUtility "example/pkg/utility"
)

var oRabbitmqCommand = &cobra.Command{
	Use:   "rabbitmq",
	Short: "啟動 rabbitmq 服務",
	Run: func(cmd *cobra.Command, args []string) {

		// 收到中斷/終止訊號時 ctx 會被取消，ConsumerRouter.Serve 監聽 ctx.Done() 後返回，
		// 不是靠 process 被系統強制殺掉才停止消費。
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitRabbitmqContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		oRabbitmqRouter := registerRabbitmq.Init(oContainer)

		pkgUtility.Logger(pkgUtility.Default).Info("啟動 RABBITMQ 服務。")

		if err := oRabbitmqRouter.Serve(ctx); err != nil {
			log.Printf("rabbitmq stopped: %v", err)
		}

	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oRabbitmqCommand)
}
