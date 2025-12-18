package auth

import (
	"go-tasks/internal/account"
	"go-tasks/internal/collaborator"
	"go-tasks/internal/pkg/jwttoken"

	"github.com/google/uuid"
)

type Service interface {
	SignIn(key, password string) (string, error)
	SignUp(account *account.Account, password string) error
	UpdatePassword(id uuid.UUID, password string) error
	LoadCollaboratorContext(collaboratorID uuid.UUID, authData jwttoken.AuthData) (string, error)
	UpdateEmail(id uuid.UUID, email string) error
	UpdatePhone(id uuid.UUID, phone string) error
}

type service struct {
	TokenService        jwttoken.TokenService
	AccountService      account.Service
	CollaboratorService collaborator.Service
}

func NewService(TokenService jwttoken.TokenService, AccountService account.Service, CollaboratorService collaborator.Service) Service {
	return &service{TokenService, AccountService, CollaboratorService}
}

func (s *service) SignIn(key, password string) (string, error) {
	account, found, err := s.AccountService.FindForSignIn(key, key, key)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrInvalidCredentials
	}
	if err = s.AccountService.VerifyPassword(password, account.Password); err != nil {
		return "", ErrInvalidCredentials
	}
	var collaboratorData jwttoken.CollaboratorData
	if len(account.Roles) > 0 {
		role := account.Roles[0]
		collaboratorData = jwttoken.CollaboratorData{
			ID:        role.ID,
			CompanyID: role.CompanyID,
			Role:      role.Role,
		}
	}
	jwt, err := s.TokenService.GenerateToken(&jwttoken.AuthData{
		Id:           account.ID,
		Username:     account.Username,
		Role:         string(account.Role),
		Collaborator: collaboratorData,
	})
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func (s *service) SignUp(account *account.Account, password string) error {
	return s.AccountService.CreateWithPassword(account, password)
}

func (s *service) UpdatePassword(id uuid.UUID, password string) error {
	return s.AccountService.UpdatePassword(id, password)
}

func (s *service) UpdateEmail(id uuid.UUID, email string) error {
	return s.AccountService.UpdateEmail(id, email)
}

func (s *service) UpdatePhone(id uuid.UUID, phone string) error {
	return s.AccountService.UpdatePhone(id, phone)
}

func (s *service) LoadCollaboratorContext(collaboratorID uuid.UUID, authData jwttoken.AuthData) (string, error) {
	collaborator, found, err := s.CollaboratorService.GetById(collaboratorID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrNotCollaborator
	}

	if found, err = s.CollaboratorService.IsAccountCollaborator(authData.Id, collaborator.ID); !found || err != nil {
		return "", ErrNotCollaborator
	}
	jwt, err := s.TokenService.GenerateToken(&jwttoken.AuthData{
		Id:       authData.Id,
		Username: authData.Username,
		Role:     authData.Role,
		Collaborator: jwttoken.CollaboratorData{
			ID:        collaborator.ID,
			CompanyID: collaborator.CompanyID,
			Role:      string(collaborator.Role),
		},
	})
	if err != nil {
		return "", err
	}
	return jwt, nil
}
