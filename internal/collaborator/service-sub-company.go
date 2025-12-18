package collaborator

import (
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceSubCompany interface {
	Create(collaborator *CollaboratorInfo, companyId uuid.UUID, role CollaboratorRole) error
	CreateWithAccount(collaborator *CompanyCollaborator, password string, companyId uuid.UUID, role CollaboratorRole) error
	Delete(id uuid.UUID, companyId uuid.UUID, role CollaboratorRole) (bool, error)
	GetByCompanyId(companyId uuid.UUID, role CollaboratorRole, p httpx.PaginationParams) ([]CompanyCollaborator, int64, error)
	GetById(id, companyId uuid.UUID, role CollaboratorRole) (*CompanyCollaborator, bool, error)
	GetSubordinateRoles(acting CollaboratorRole) []CollaboratorRole
	Update(collaborator *CompanyCollaborator, companyId uuid.UUID, role CollaboratorRole) error
}

type serviceSubCompany struct {
	repository     Repository
	serviceCompany ServiceCompany
}

func NewServiceSubCompany(r Repository, a ServiceCompany) ServiceSubCompany {
	return &serviceSubCompany{repository: r, serviceCompany: a}
}

func (s *serviceSubCompany) GetByCompanyId(companyId uuid.UUID, role CollaboratorRole, p httpx.PaginationParams) ([]CompanyCollaborator, int64, error) {
	roles := s.GetSubordinateRoles(role)
	return s.repository.GetByCompanyIdAndRoles(p, companyId, roles)
}

func (s *serviceSubCompany) GetById(id, companyId uuid.UUID, role CollaboratorRole) (*CompanyCollaborator, bool, error) {
	item, ok, err := s.serviceCompany.GetById(id, companyId)
	if ok && err == nil {
		if innerErr := s.VerifySubordinate(role, item.Role); innerErr != nil {
			return item, ok, innerErr
		}
	}
	return item, ok, err
}
func (s *serviceSubCompany) Update(collaborator *CompanyCollaborator, companyId uuid.UUID, acting CollaboratorRole) error {
	if err := s.VerifyCollaboratorSubordinate(collaborator.ID, acting); err != nil {
		return err
	}
	if err := s.VerifySubordinate(acting, collaborator.Role); err != nil {
		return err
	}
	return s.serviceCompany.Update(collaborator, companyId)
}

func (s *serviceSubCompany) Create(collaborator *CollaboratorInfo, companyId uuid.UUID, role CollaboratorRole) error {
	if err := s.VerifySubordinate(role, collaborator.Role); err != nil {
		return err
	}
	return s.serviceCompany.Create(collaborator, companyId)

}
func (s *serviceSubCompany) CreateWithAccount(collaborator *CompanyCollaborator, password string, companyId uuid.UUID, role CollaboratorRole) error {
	if err := s.VerifySubordinate(role, collaborator.Role); err != nil {
		return err
	}
	return s.serviceCompany.CreateWithAccount(collaborator, password, companyId)
}
func (s *serviceSubCompany) Delete(id uuid.UUID, companyId uuid.UUID, acting CollaboratorRole) (bool, error) {
	if err := s.VerifyCollaboratorSubordinate(id, acting); err != nil {
		return false, err
	}
	return s.serviceCompany.Delete(id, companyId)
}

func (s *serviceSubCompany) VerifyCollaboratorSubordinate(id uuid.UUID, acting CollaboratorRole) error {
	role, ok, err := s.repository.GetRoleById(id)
	if err != nil {
		return err
	}
	if ok {
		if err := s.VerifySubordinate(acting, role); err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceSubCompany) VerifySubordinate(acting CollaboratorRole, subordinate CollaboratorRole) error {
	for _, role := range s.GetSubordinateRoles(acting) {
		if role == subordinate {
			return nil
		}
	}
	return ErrForbiddenNotSubordinate
}
func (s *serviceSubCompany) GetSubordinateRoles(acting CollaboratorRole) []CollaboratorRole {
	switch acting {
	case CollaboratorRoleOwner:
		return []CollaboratorRole{CollaboratorRoleManager, CollaboratorRoleOps}
	case CollaboratorRoleManager:
		return []CollaboratorRole{CollaboratorRoleOps}
	default:
		return nil
	}
}
