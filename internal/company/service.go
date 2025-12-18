package company

import (
	"go-tasks/internal/account"
	"go-tasks/internal/collaborator"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type Service interface {
	Create(company *CompanyWithOwner) (Company, error)
	Delete(id uuid.UUID) (bool, error)
	Get(p httpx.PaginationParams) ([]Company, int64, error)
	GetById(id uuid.UUID) (*Company, bool, error)
	Update(company *Company) error
}

type service struct {
	repository             Repository
	collaboratorRepository collaborator.Repository
	accountService         account.Service
}

func NewService(r Repository, collabRepo collaborator.Repository, accountService account.Service) Service {
	return &service{repository: r, collaboratorRepository: collabRepo, accountService: accountService}
}

func (s *service) Create(company *CompanyWithOwner) (Company, error) {
	var created Company
	tx := s.repository.GetConnection().BeginTransaction()
	// 1) Create company
	comp := Company{CompanyInfo: CompanyInfo{Title: company.Title}}
	if err := s.repository.CreateTx(tx, &comp); err != nil {
		tx.Rollback()
		return created, err
	}

	// 2) Create owner account via account service (handles hashing)
	acc := account.Account{
		Username: company.Owner.Account.Username,
		Name:     company.Owner.Account.Name,
		Phone:    company.Owner.Account.Phone,
		Email:    company.Owner.Account.Email,
		Role:     account.RoleUser, // owner as USER
	}
	if err := s.accountService.CreateWithPasswordTx(tx, &acc, company.Owner.Account.Password); err != nil {
		tx.Rollback()
		return created, err
	}

	// 3) Link owner as collaborator to the company
	collaborator := collaborator.CollaboratorInfo{
		CompanyID: comp.ID,
		AccountID: acc.ID,
		Role:      collaborator.CollaboratorRoleOwner,
	}
	if err := s.collaboratorRepository.CreateTx(tx, &collaborator); err != nil {
		tx.Rollback()
		return created, err
	}

	created = comp
	if err := tx.Commit(); err != nil {
		return created, err
	}
	return created, nil
}

func (s *service) Delete(id uuid.UUID) (bool, error) {
	return s.repository.Delete(id)
}

func (s *service) Get(p httpx.PaginationParams) ([]Company, int64, error) {
	return s.repository.Get(p)
}

func (s *service) GetById(id uuid.UUID) (*Company, bool, error) {
	return s.repository.GetById(id)
}

func (s *service) Update(company *Company) error {
	return s.repository.Update(company)
}
