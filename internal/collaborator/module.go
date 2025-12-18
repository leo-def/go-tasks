package collaborator

import (
	"go-tasks/internal/account"
	"go-tasks/internal/pkg/database"
)

type Module struct {
	Controller           *Controller
	ControllerCompany    *ControllerCompany
	ControllerSubCompany *ControllerSubCompany
	Service              *Service
	ServiceCompany       *ServiceCompany
	ServiceSubCompany    *ServiceSubCompany
}

func Initialize(db database.Connection, accountService account.Service) *Module {
	repository := NewRepository(db)
	service := NewService(repository, accountService)
	serviceCompany := NewServiceCompany(repository, service, accountService)
	serviceSubCompany := NewServiceSubCompany(repository, serviceCompany)
	controller := NewController(service)
	controllerCompany := NewControllerCompany(serviceCompany)
	controllerSubCompany := NewControllerSubCompany(serviceSubCompany)
	return &Module{
		Controller:           controller,
		ControllerCompany:    controllerCompany,
		ControllerSubCompany: controllerSubCompany,
		Service:              &service,
		ServiceCompany:       &serviceCompany,
		ServiceSubCompany:    &serviceSubCompany,
	}
}
