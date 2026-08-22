package register

import (
	pkg "example/pkg"

	container "example/container"
)

func RabbitmqInit(oContainer *container.RabbitmqContainer) *pkg.ConsumerRouter {
	oRouter := pkg.NewConsumerRouter(oContainer.Conn)
	oRouter.HandleFunc("Admin.Resource.AppUser.IncreaseBalance", oContainer.ConsumerAdminResourceAppUser.IncreaseBalance)

	return oRouter
}
