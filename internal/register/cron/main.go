package registerCron

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	container "example/container"
	pkgUtility "example/pkg/utility"
)

func Init(oContainer *container.CronContainer) *cron.Cron {
	oCron := cron.New()

	// 這裡 oContainer.CronAdminResourceAppUser.IncreaseBalance 是，閉包，還沒執行，所以啟動不會連 mysql
	if _, err := oCron.AddFunc("* * * * *", oContainer.CronAdminResourceAppUser.IncreaseBalance); err != nil {
		oLogger := pkgUtility.Logger(pkgUtility.Cron)
		oLogger.Fatal("cron: failed to register CronAppUser.IncreaseBalance job", zap.Error(err))
	}

	if _, err := oCron.AddFunc("* * * * *", oContainer.CronAdminAuthenticationAuthenticator.SignIn); err != nil {
		oLogger := pkgUtility.Logger(pkgUtility.Cron)
		oLogger.Fatal("cron: failed to register CronAdminAuthenticationAuthenticator.SignIn job", zap.Error(err))
	}

	return oCron
}
