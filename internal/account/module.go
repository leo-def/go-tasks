package account

import (
    "go-tasks/internal/pkg/database"
)

type Module struct {
	Service *Service
}

func Initialize(db database.Connection) *Module {
    repository := NewRepository(db)
    service := NewService(repository)
    return &Module{Service: &service}
}
