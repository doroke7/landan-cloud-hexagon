package registerRabbitmq

import (
	pkg "example/pkg"

	container "example/container"
)

func Init(oContainer *container.RabbitmqContainer) *pkg.RabbitmqRouter {
	oRouter := pkg.NewRabbitmqRouter(oContainer.Conn)
	oRouter.HandleFunc("Admin.Resource.AppUser.IncreaseBalance", oContainer.ConsumerAdminResourceAppUser.IncreaseBalance)

	return oRouter
}
