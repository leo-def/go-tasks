package lifecycle

import (
    "go-tasks/internal/pkg/database"
)

type Module struct {
	Service Service
}

func Initialize(db database.Connection) *Module {
    repo := NewRepository(db)
    validator := NewValidationService()
    service := NewService(repo, validator)
    return &Module{Service: service}
}
