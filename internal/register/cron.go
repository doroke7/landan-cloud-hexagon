package register

import (
	"log"

	"github.com/robfig/cron/v3"

	container "example/container"
)

func CronInit(oContainer *container.CronContainer) *cron.Cron {
	oCron := cron.New()

	// 這裡 oContainer.CronAdminResourceAppUser.IncreaseBalance 是，閉包，還沒執行，所以啟動不會連 mysql
	if _, err := oCron.AddFunc("* * * * *", oContainer.CronAdminResourceAppUser.IncreaseBalance); err != nil {
		log.Fatalf("cron: failed to register CronAppUser.IncreaseBalance job: %v", err)
	}

	if _, err := oCron.AddFunc("* * * * *", oContainer.CronAdminAuthenticationAuthenticator.SignIn); err != nil {
		log.Fatalf("cron: failed to register CronAdminAuthenticationAuthenticator.SignIn job: %v", err)
	}

	return oCron
}
