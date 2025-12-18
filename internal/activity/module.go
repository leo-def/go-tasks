package activity

import (
    "go-tasks/internal/lifecycle"
    "go-tasks/internal/pkg/database"
)

type Module struct {
	Controller     *ControllerCompany
	ControllerOwn  *ControllerOwn
	Service        *Service
	ServiceCompany *ServiceCompany
	ServiceOwn     *ServiceOwn
}

func Initialize(db database.Connection, lifecycleService lifecycle.Service) *Module {
    repo := NewRepository(db)
    service := NewService(repo, lifecycleService)
    serviceCompany := NewServiceCompany(repo, service)
    serviceOwn := NewServiceOwn(repo, service)
    controller := NewControllerCompany(serviceCompany)
    controllerOwn := NewControllerOwn(serviceOwn)
    return &Module{Controller: controller, ControllerOwn: controllerOwn, Service: &service, ServiceCompany: &serviceCompany, ServiceOwn: &serviceOwn}
}
