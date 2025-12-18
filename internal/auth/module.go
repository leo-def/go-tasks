package auth

import (
	"go-tasks/internal/account"
	"go-tasks/internal/collaborator"
	"go-tasks/internal/pkg/jwttoken"
)

type Module struct {
	Controller   *Controller
	TokenService jwttoken.TokenService
}

func Initialize(accountService account.Service, collaboratorService collaborator.Service) *Module {
	tokenService := jwttoken.NewService()
	service := NewService(tokenService, accountService, collaboratorService)
	controller := NewController(service)
	return &Module{Controller: controller, TokenService: tokenService}
}
