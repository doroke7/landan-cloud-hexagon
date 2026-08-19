package cmd

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	pkg "example/pkg"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
)

var oResourceCommand = &cobra.Command{
	Use:   "resource",
	Short: "啟動 Resource 服務",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitResourceContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}

		oResourceServer := register.ResourceInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.RESOURCE.PORT)
		if err != nil {
			log.Fatal(err)
		}
		pkg.Logger(pkg.Default).Info("啟動 RESOURCE 服務。 port: " + bootstrap.CONFIG.SERVICES.RESOURCE.PORT)

		// 收到中斷/終止訊號時 ctx 會被取消，主動 GracefulStop，讓 gRPC 停止 accept 新連線、
		// 關掉 listener，Serve() 才會正常返回並釋放 port，
		// 不是靠 process 被系統強制殺掉才釋放。
		go func() {
			<-ctx.Done()
			oResourceServer.GracefulStop()
		}()

		if err := oResourceServer.Serve(oListener); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oResourceCommand)
}
