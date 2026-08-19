package register

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
	container "example/container"
)

// CommandInit 只組裝子命令的「形狀」（名字、flag），不在這裡連任何基礎設施；
// container.InitCommandContainer() 延後到 Run: 真正執行時才呼叫，
// 這樣註冊子命令樹（init() 階段）就不需要先連上 MySQL。
func CommandInit(oCommandCommand *cobra.Command) *cobra.Command {
	var iId uint
	var iAmount uint

	oAppUserIncreaseBalanceCommand := &cobra.Command{
		Use:   "Admin-Resource-AppUser-IncreaseBalance",
		Short: "AppUser-IncreaseBalance 相關命令",
		Run: func(oCmd *cobra.Command, args []string) {
			oContainer, err := container.InitCommandContainer(context.Background())
			if err != nil {
				log.Fatalf("command: failed to init container: %v", err)
			}

			if err := oContainer.CommandAdminReourceAppUser.IncreaseBalance(iId, iAmount); err != nil {
				log.Printf("increase balance failed: %v", err)
			}
		},
	}

	oAppUserIncreaseBalanceCommand.Flags().UintVar(&iId, "id", 1, "AppUser 的 id")
	oAppUserIncreaseBalanceCommand.Flags().UintVar(&iAmount, "amount", 10, "要增加的餘額")

	oCommandCommand.AddCommand(oAppUserIncreaseBalanceCommand)

	var sName string
	var sPassword string

	oAuthenticatorSignInCommand := &cobra.Command{
		Use:   "Admin-Authentication-Authenticator-SignIn",
		Short: "Authenticator-SignIn 相關命令",
		Run: func(oCmd *cobra.Command, args []string) {
			oContainer, err := container.InitCommandContainer(context.Background())
			if err != nil {
				log.Fatalf("command: failed to init container: %v", err)
			}

			// NOTE: command carrier 目前沒有自己的 services.command.admin.jwt 設定，先借用 http 那組 secret。
			sAuthorization, err := oContainer.CommandAdminAuthenticationSignIn.SignIn(sName, sPassword, bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET)
			if err != nil {
				log.Printf("sign in failed: %v", err)
				return
			}

			log.Printf("sign in succeeded, authorization: %s", sAuthorization)
		},
	}

	oAuthenticatorSignInCommand.Flags().StringVar(&sName, "name", "", "登入帳號")
	oAuthenticatorSignInCommand.Flags().StringVar(&sPassword, "password", "", "登入密碼")

	oCommandCommand.AddCommand(oAuthenticatorSignInCommand)

	return oCommandCommand
}
