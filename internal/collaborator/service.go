package collaborator

import (
	"go-tasks/internal/account"
	"go-tasks/internal/pkg/token"
	"time"

	"github.com/google/uuid"
)

type Service interface {
    Create(collaborator *CollaboratorInfo) error
    CreateWithAccount(collaborator *Collaborator, password string) error
    GetByAccountId(accountId uuid.UUID) ([]AccountCollaborator, int64, error)
    GetById(id uuid.UUID) (Collaborator, bool, error)
    IsAccountCollaborator(accountId, collaboratorId uuid.UUID) (bool, error)
    Update(collaborator *Collaborator) error
}

type service struct {
	repository     Repository
	accountService account.Service
}

func NewService(r Repository, a account.Service) Service {
	return &service{r, a}
}

func (s *service) GetByAccountId(accountId uuid.UUID) ([]AccountCollaborator, int64, error) {
	return s.repository.GetByAccountId(accountId)
}

func (s *service) GetById(id uuid.UUID) (Collaborator, bool, error) {
	return s.repository.GetById(id)
}

func (s *service) Delete(id uuid.UUID) (bool, error) {
	return s.repository.Delete(id)
}

func (s *service) Update(collaborator *Collaborator) error {
	return s.repository.Update(collaborator)
}

func (s *service) Create(collaborator *CollaboratorInfo) error {
	tx := s.repository.GetConnection().BeginTransaction()
	collaborator.Active = false
	if err := s.repository.CreateTx(tx, collaborator); err != nil {
		tx.Rollback()
		return err
	}
	tok := token.GenerateActivationTokenFrom(collaborator.ID, time.Now())
	if err := s.repository.SetActivationTokenTx(tx, collaborator.ID, tok); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *service) CreateWithAccount(collaborator *Collaborator, password string) error {
	tx := s.repository.GetConnection().BeginTransaction()
	accountModel := &account.Account{
		Username: collaborator.Account.Username,
		Name:     collaborator.Account.Name,
		Phone:    collaborator.Account.Phone,
		Email:    collaborator.Account.Email,
		Role:     account.RoleUser,
	}
	if err := s.accountService.CreateWithPasswordTx(tx, accountModel, password); err != nil {
		tx.Rollback()
		return err
	}
	collaborator.AccountID = accountModel.ID
	ci := &CollaboratorInfo{
		CompanyID: collaborator.CompanyID,
		AccountID: collaborator.AccountID,
		Role:      collaborator.Role,
		Active:    true,
	}
	if err := s.repository.CreateTx(tx, ci); err != nil {
		tx.Rollback()
		return err
	}
	collaborator.ID = ci.ID
	return tx.Commit()
}

func (s *service) IsAccountCollaborator(accountId uuid.UUID, collaboratorId uuid.UUID) (bool, error) {
	_, found, err := s.repository.GetByIdAndAccountId(collaboratorId, accountId)
	return found, err
}
