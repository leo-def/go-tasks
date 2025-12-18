package collaborator

import (
	"go-tasks/internal/account"
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceCompany interface {
	Update(collaborator *CompanyCollaborator, companyId uuid.UUID) error
	Create(collaborator *CollaboratorInfo, companyId uuid.UUID) error
	CreateWithAccount(collaborator *CompanyCollaborator, password string, companyId uuid.UUID) error
	Get(p httpx.PaginationParams, companyId uuid.UUID) ([]CompanyCollaborator, int64, error)
	GetById(id, companyId uuid.UUID) (*CompanyCollaborator, bool, error)
	Delete(id uuid.UUID, companyId uuid.UUID) (bool, error)
	IsAccountCompanyCollaborator(accountId, companyId uuid.UUID) (bool, error)
}

type serviceCompany struct {
	repository          Repository
	collaboratorService Service
	accountService      account.Service
}

func NewServiceCompany(repository Repository, collaboratorService Service, accountService account.Service) ServiceCompany {
	return &serviceCompany{repository: repository, collaboratorService: collaboratorService, accountService: accountService}
}

func (s *serviceCompany) Update(collaborator *CompanyCollaborator, companyId uuid.UUID) error {
	err := s.VerifyCompanyId(collaborator.ID, companyId)
	if err != nil {
		return err
	}
	param := Collaborator{
		CollaboratorInfo: collaborator.CollaboratorInfo,
		Account:          collaborator.Account,
	}
	if err = s.collaboratorService.Update(&param); err == nil {
		collaborator.ID = param.ID
	}
	return err
}

func (s *serviceCompany) Create(collaborator *CollaboratorInfo, companyId uuid.UUID) error {
	collaborator.CompanyID = companyId
	return s.collaboratorService.Create(collaborator)
}

func (s *serviceCompany) CreateWithAccount(collaborator *CompanyCollaborator, password string, companyId uuid.UUID) error {
	collaborator.CompanyID = companyId
	param := Collaborator{
		CollaboratorInfo: collaborator.CollaboratorInfo,
		Account:          collaborator.Account,
	}
	err := s.collaboratorService.Update(&param)
	if err == nil {
		collaborator.ID = param.ID
	}
	err = s.collaboratorService.CreateWithAccount(&param, password)
	if err == nil {
		collaborator.ID = param.ID
	}
	return err
}

func (s *serviceCompany) Get(p httpx.PaginationParams, companyId uuid.UUID) ([]CompanyCollaborator, int64, error) {
	return s.repository.GetByCompanyId(p, companyId)
}

func (s *serviceCompany) GetById(id uuid.UUID, companyId uuid.UUID) (*CompanyCollaborator, bool, error) {
	return s.repository.GetCompanyCollaboratorById(id, companyId)
}

func (s *serviceCompany) Delete(id uuid.UUID, companyId uuid.UUID) (bool, error) {
	err := s.VerifyCompanyId(id, companyId)
	if err != nil {
		return false, err
	}
	info, found, err := s.repository.GetInfoById(id)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	tx := s.repository.GetConnection().BeginTransaction()
	txdb, _ := database.AsGormTx(tx)
	if err = txdb.Exec("DELETE FROM assignments WHERE assigner_id = ?", id).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	ok, err := s.repository.DeleteTx(tx, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return false, err
	}
	_, count, err := s.repository.GetByAccountId(info.AccountID)
	if err != nil {
		return ok, err
	}
	if count == 0 {
		if acc, exists, err := s.accountService.GetById(info.AccountID); err == nil && exists && acc.Role == account.RoleUser {
			_, _ = s.accountService.Delete(info.AccountID)
		}
	}
	return ok, nil
}

func (s *serviceCompany) IsAccountCompanyCollaborator(accountId, companyId uuid.UUID) (bool, error) {
	_, found, err := s.repository.GetByAccountIdAndCompanyId(accountId, companyId)
	return found, err
}

func (s *serviceCompany) VerifyCompanyId(id, companyId uuid.UUID) error {
	found, err := s.repository.CheckByIdAndCompanyId(id, companyId)
	if err != nil {
		return err
	}
	if !found {
		return ErrForbiddenNotInCompany
	}
	return nil
}
