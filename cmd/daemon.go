package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	container "example/container"
	registerDaemon "example/internal/register/daemon"
	pkgUtility "example/pkg/utility"
)

var oDaemonCommand = &cobra.Command{
	Use:   "daemon",
	Short: "啟動 daemon 服務",
	Run: func(cmd *cobra.Command, args []string) {
		oCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oLogger := pkgUtility.Logger(pkgUtility.Deamon)

		oContainer, err := container.InitDaemonContainer(oCtx)
		if err != nil {
			oErrorField := zap.Error(err)
			oLogger.Fatal("初始化 daemon container 失敗", oErrorField)
		}

		oDaemonRouter := registerDaemon.Init(oContainer)

		oLogger.Info("啟動 DAEMON 服務。")

		if err := oDaemonRouter.Serve(oCtx); err != nil {
			oErrorField := zap.Error(err)
			oLogger.Error("DAEMON 服務已停止", oErrorField)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oDaemonCommand)
}
