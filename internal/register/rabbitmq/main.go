package registerRabbitmq

import (
	pkgRabbitmq "example/pkg/rabbitmq"

	container "example/container"
)

func Init(oContainer *container.RabbitmqContainer) *pkgRabbitmq.RabbitmqRouter {
	oRouter := pkgRabbitmq.NewRabbitmqRouter(oContainer.Conn)
	oRouter.HandleFunc("Admin.Resource.AppUser.IncreaseBalance", oContainer.ConsumerAdminResourceAppUser.IncreaseBalance)

	return oRouter
}
