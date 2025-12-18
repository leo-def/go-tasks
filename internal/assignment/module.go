package assignment

import (
	"go-tasks/internal/participation"
	"go-tasks/internal/pkg/database"
)

type Module struct {
	Controller *ControllerTask
	Service    *ServiceTask
}

func Initialize(db database.Connection, participationService participation.Service) *Module {
	repository := NewRepository(db)
	service := NewServiceTask(repository, participationService)
	controller := NewControllerTask(service)
	return &Module{Controller: controller, Service: &service}
}
