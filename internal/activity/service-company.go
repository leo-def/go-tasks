package activity

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceCompany interface {
	Create(activity *CompanyActivity, createdBy uuid.UUID, companyID uuid.UUID) error
	Delete(id, companyID uuid.UUID) (bool, error)
	Get(p httpx.PaginationParams, companyID uuid.UUID) ([]CompanyActivity, int64, error)
	GetById(id, companyID uuid.UUID) (*CompanyActivity, bool, error)
	Update(activity *CompanyActivity, companyID uuid.UUID) error
	UpdatedStatus(id uuid.UUID, status lifecycle.LifecycleStatus, companyID uuid.UUID) error
}

type serviceCompany struct {
	repository Repository
	service    Service
}

func NewServiceCompany(repository Repository, service Service) ServiceCompany {
	return &serviceCompany{repository, service}
}

func (s *serviceCompany) Get(p httpx.PaginationParams, companyID uuid.UUID) ([]CompanyActivity, int64, error) {
	return s.repository.GetByCompanyId(p, companyID)
}

func (s *serviceCompany) GetById(id, companyID uuid.UUID) (*CompanyActivity, bool, error) {
	return s.repository.GetCompanyActivityById(id, companyID)
}

func (s *serviceCompany) Delete(id, companyID uuid.UUID) (bool, error) {
	if err := s.VerifyCompanyId(id, companyID); err != nil {
		return false, err
	}
	return s.service.Delete(id)
}

func (s *serviceCompany) Update(activity *CompanyActivity, companyID uuid.UUID) error {
	if err := s.VerifyCompanyId(activity.ID, companyID); err != nil {
		return err
	}
	param := Activity{
		ActivityInfo: activity.ActivityInfo,
		Company: ActivityCompany{
			ID: companyID,
		},
		Lifecycle: activity.Lifecycle,
		Owner:     activity.Owner,
		CreatedBy: activity.CreatedBy,
	}
	if err := s.service.Update(&param); err != nil {
		return err
	}
	activity.ID = param.ID
	if updated, found, err := s.repository.GetCompanyActivityById(activity.ID, companyID); err != nil {
		return err
	} else if found {
		*activity = *updated
	}
	return nil
}

func (s *serviceCompany) Create(activity *CompanyActivity, createdBy uuid.UUID, companyID uuid.UUID) error {
	activity.CompanyID = companyID
	param := Activity{
		ActivityInfo: activity.ActivityInfo,
		Company: ActivityCompany{
			ID: companyID,
		},
		Lifecycle: activity.Lifecycle,
		Owner:     activity.Owner,
		CreatedBy: activity.CreatedBy,
	}
	if err := s.service.Create(&param, createdBy); err != nil {
		return err
	}
	activity.ID = param.ID
	if created, found, err := s.repository.GetCompanyActivityById(activity.ID, companyID); err != nil {
		return err
	} else if found {
		*activity = *created
	}
	return nil
}

func (s *serviceCompany) UpdatedStatus(id uuid.UUID, status lifecycle.LifecycleStatus, companyID uuid.UUID) error {
	if err := s.VerifyCompanyId(id, companyID); err != nil {
		return err
	}
	return s.service.UpdatedStatus(id, status)
}

func (s *serviceCompany) VerifyCompanyId(id uuid.UUID, companyID uuid.UUID) error {
	found, err := s.repository.CheckByIdAndCompanyId(id, companyID)
	if err != nil {
		return err
	}
	if !found {
		return ErrForbiddenNotInCompany
	}
	return nil
}
