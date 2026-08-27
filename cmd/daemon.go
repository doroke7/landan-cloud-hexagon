package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	registerDaemon "example/internal/register/daemon"
	pkgUtility "example/pkg/utility"
)

var oDaemonCommand = &cobra.Command{
	Use:   "daemon",
	Short: "啟動 daemon 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, _ := container.InitDaemonContainer(ctx)

		oDaemonRouter := registerDaemon.Init(oContainer)

		pkgUtility.Logger(pkgUtility.Default).Info("啟動 DAEMON 服務。")

		oDaemonRouter.Serve(ctx)
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oDaemonCommand)
}
