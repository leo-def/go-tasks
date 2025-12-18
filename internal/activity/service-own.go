package activity

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceOwn interface {
	Create(activity *OwnActivity, createdBy uuid.UUID, ownerID uuid.UUID) error
	Delete(id, ownerID uuid.UUID) (bool, error)
	Get(p httpx.PaginationParams, ownerID uuid.UUID) ([]OwnActivity, int64, error)
	GetById(id, ownerID uuid.UUID) (*OwnActivity, bool, error)
	Update(activity *OwnActivity, ownerID uuid.UUID) error
	UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, ownerID uuid.UUID) error
}

type serviceOwn struct {
	repository Repository
	service    Service
}

func NewServiceOwn(repository Repository, service Service) ServiceOwn {
	return &serviceOwn{repository, service}
}

func (s *serviceOwn) Get(p httpx.PaginationParams, ownerID uuid.UUID) ([]OwnActivity, int64, error) {
	return s.repository.GetByOwnerId(p, ownerID)
}

func (s *serviceOwn) GetById(id uuid.UUID, ownerID uuid.UUID) (*OwnActivity, bool, error) {
	if err := s.VerifyOwnership(id, ownerID); err != nil {
		return nil, false, err
	}
	activity, found, err := s.repository.GetOwnActivityById(id, ownerID)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	return activity, true, nil
}

func (s *serviceOwn) Delete(id uuid.UUID, ownerID uuid.UUID) (bool, error) {
	err := s.VerifyOwnership(id, ownerID)
	if err != nil {
		return false, err
	}
	return s.service.Delete(id)
}

func (s *serviceOwn) Update(activity *OwnActivity, ownerID uuid.UUID) error {
	if err := s.VerifyOwnership(activity.ID, ownerID); err != nil {
		return err
	}
	param := Activity{
		ActivityInfo: activity.ActivityInfo,
		Lifecycle:    activity.Lifecycle,
		CreatedBy:    activity.CreatedBy,
	}
	if err := s.service.Update(&param); err != nil {
		return err
	}
	activity.ID = param.ID
	if updated, found, err := s.repository.GetOwnActivityById(activity.ID, ownerID); err != nil {
		return err
	} else if found {
		*activity = *updated
	}
	return nil
}

func (s *serviceOwn) Create(activity *OwnActivity, createdBy uuid.UUID, ownerID uuid.UUID) error {
	activity.OwnerID = ownerID
	param := Activity{
		ActivityInfo: activity.ActivityInfo,
		Lifecycle:    activity.Lifecycle,
		CreatedBy:    activity.CreatedBy,
	}
	if err := s.service.Create(&param, createdBy); err != nil {
		return err
	}
	activity.ID = param.ID
	if created, found, err := s.repository.GetOwnActivityById(activity.ID, ownerID); err != nil {
		return err
	} else if found {
		*activity = *created
	}
	return nil
}

func (s *serviceOwn) UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, ownerID uuid.UUID) error {
	err := s.VerifyOwnership(id, ownerID)
	if err != nil {
		return err
	}
	return s.service.UpdatedStatus(id, status)
}

func (s *serviceOwn) VerifyOwnership(id uuid.UUID, ownerID uuid.UUID) error {
	found, err := s.repository.CheckByIdAndOwnerId(id, ownerID)
	if err != nil {
		return err
	}
	if !found {
		return ErrForbiddenNotOwner
	}
	return nil
}
