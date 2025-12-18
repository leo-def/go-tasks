package rating

import (
    "go-tasks/internal/participation"
    "go-tasks/internal/pkg/database"
)

type Module struct {
	Controller *Controller
}

func Initialize(db database.Connection, participationService participation.Service) *Module {
    repo := NewRepository(db)
    service := NewService(repo, participationService)
    controller := NewController(service)
    return &Module{Controller: controller}
}
