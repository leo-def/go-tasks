package company

import (
    "go-tasks/internal/account"
    "go-tasks/internal/collaborator"
    "go-tasks/internal/pkg/database"
)

type Module struct {
	Controller *Controller
}

// Initialize wires repository and service, injecting account service for transactional account creation.
func Initialize(db database.Connection, accountService account.Service) *Module {
    repo := NewRepository(db)
    collabRepo := collaborator.NewRepository(db)
    service := NewService(repo, collabRepo, accountService)
    controller := NewController(service)

    return &Module{Controller: controller}
}
