package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	container "example/container"
	register "example/internal/register"
)

var oCronCommand = &cobra.Command{
	Use:   "cron",
	Short: "啟動排程服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, _ := container.InitCronContainer(ctx)

		oCron := register.CronInit(oContainer)
		oCron.Start()

		<-ctx.Done()
		oCron.Stop()
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oCronCommand)
}
