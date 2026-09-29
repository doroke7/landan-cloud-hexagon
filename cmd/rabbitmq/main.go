package rabbitmq

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	container "example/container"
	registerRabbitmq "example/internal/register/rabbitmq"
	pkgUtility "example/pkg/utility"
)

var Command = &cobra.Command{
	Use:   "rabbitmq",
	Short: "啟動 rabbitmq 服務",
	Run: func(cmd *cobra.Command, args []string) {

		// 收到中斷/終止訊號時 oCtx 會被取消，ConsumerRouter.Serve 監聽 oCtx.Done() 後返回，
		// 不是靠 process 被系統強制殺掉才停止消費。
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Consumer)

		oContainer, err := container.InitRabbitmqContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 rabbitmq container 失敗", oErrorField)
		}

		oRabbitmqRouter := registerRabbitmq.Init(oContainer)

		oLogger.Info("啟動 RABBITMQ 服務。")

		if err := oRabbitmqRouter.Serve(oCtx); err != nil {
			oErrorField := zap.Error(err)
			oLogger.Error("RABBITMQ 服務已停止", oErrorField)
		}
	},
}
