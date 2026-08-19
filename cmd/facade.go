package cmd

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
	register "example/internal/register"
	pkg "example/pkg"
)

var oFacadeCommand = &cobra.Command{
	Use:   "facade",
	Short: "啟動 Facade 服務",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		oContainer, err := container.InitFacadeContainer(ctx)
		if err != nil {
			log.Fatal(err)
		}
		oFacadeServer := register.FacadeInit(oContainer)

		oListener, err := net.Listen("tcp", ":"+bootstrap.CONFIG.SERVICES.FACADE.PORT)
		if err != nil {
			log.Fatal(err)
		}
		pkg.Logger(pkg.Default).Info("啟動 FACADE 服務。 port: " + bootstrap.CONFIG.SERVICES.FACADE.PORT)

		go func() {
			<-ctx.Done()
			oFacadeServer.GracefulStop()
		}()

		if err := oFacadeServer.Serve(oListener); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oFacadeCommand)
}
