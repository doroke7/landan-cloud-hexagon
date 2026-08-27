package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	registerCron "example/internal/register/cron"
	pkgUtility "example/pkg/utility"
)

var oCronCommand = &cobra.Command{
	Use:   "cron",
	Short: "啟動排程服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, _ := container.InitCronContainer(ctx)

		oCron := registerCron.Init(oContainer)

		pkgUtility.Logger(pkgUtility.Default).Info("啟動 CRON 服務。")

		oCron.Start()

		<-ctx.Done()
		oCron.Stop()
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCronCommand)
}
