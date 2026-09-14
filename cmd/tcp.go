package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	container "example/container"
	registerTcp "example/internal/register/tcp"
	pkgUtility "example/pkg/utility"
)

var oTcpCommand = &cobra.Command{
	Use:   "tcp",
	Short: "啟動 TCP 服務",
	Run: func(cmd *cobra.Command, args []string) {
		// 收到中斷/終止訊號時 oCtx 會被取消，Tcp.Serve 內部監聽 oCtx.Done() 自己關掉 listener，
		// 不是靠 process 被系統強制殺掉才釋放 port。
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Tcp)

		oContainer, err := container.InitTcpContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 tcp container 失敗", oErrorField)
		}

		oTcpRouter := registerTcp.Init(oContainer)

		oLogger.Info("啟動 TCP 服務。 port: " + bootstrap.CONFIG.SERVICES.TCP.PORT)

		sAddress := ":" + bootstrap.CONFIG.SERVICES.TCP.PORT
		if err := oTcpRouter.Serve(oCtx, sAddress); err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("TCP server 異常結束", oErrorField)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oTcpCommand)
}
