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

		oContainer, err := container.InitDaemonContainer(oCtx)
		if err != nil {
			pkgUtility.Logger(pkgUtility.Deamon).Fatal("初始化 daemon container 失敗", zap.Error(err))
		}

		oDaemonRouter := registerDaemon.Init(oContainer)

		pkgUtility.Logger(pkgUtility.Deamon).Info("啟動 DAEMON 服務。")

		if err := oDaemonRouter.Serve(oCtx); err != nil {
			pkgUtility.Logger(pkgUtility.Deamon).Error("DAEMON 服務已停止", zap.Error(err))
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oDaemonCommand)
}
