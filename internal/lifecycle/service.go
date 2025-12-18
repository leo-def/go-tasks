package lifecycle

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
)

type Service interface {
    Create(lifecycleInfo *LifecycleInfo) error
    CreateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error
    Delete(uuid uuid.UUID) error
    DeleteTx(tx database.Transaction, uuid uuid.UUID) error
    Update(lifecycleInfo *LifecycleInfo) error
    UpdateStatus(id uuid.UUID, status LifecycleStatus) error
    UpdateStatusTx(tx database.Transaction, id uuid.UUID, status LifecycleStatus) error
    UpdateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error
}

type service struct {
	repository Repository
	validator  LifecycleValidation
}

func NewService(repository Repository, validator LifecycleValidation) Service {
	return &service{repository, validator}
}

func (s *service) Create(lifecycleInfo *LifecycleInfo) error {
	err := s.beforeCreate(lifecycleInfo)
	if err != nil {
		return err
	}
	return s.repository.Create(lifecycleInfo)
}

func (s *service) CreateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error {
	err := s.beforeCreate(lifecycleInfo)
	if err != nil {
		return err
	}
	return s.repository.CreateTx(tx, lifecycleInfo)
}

func (s *service) Update(lifecycleInfo *LifecycleInfo) error {
	err := s.beforeUpdate(lifecycleInfo)
	if err != nil {
		return err
	}
	return s.repository.Update(lifecycleInfo)
}

func (s *service) UpdateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error {
	err := s.beforeUpdate(lifecycleInfo)
	if err != nil {
		return err
	}
	return s.repository.UpdateTx(tx, lifecycleInfo)
}

func (s *service) Delete(uuid uuid.UUID) error {
	if err := s.beforeDelete(uuid); err != nil {
		return err
	}
	return s.repository.Delete(uuid)
}

func (s *service) DeleteTx(tx database.Transaction, uuid uuid.UUID) error {
	if err := s.beforeDelete(uuid); err != nil {
		return err
	}
	return s.repository.DeleteTx(tx, uuid)
}

func (s *service) UpdateStatus(id uuid.UUID, status LifecycleStatus) error {
	old, err := s.beforeUpdateStatus(id, status)
	if err != nil {
		return err
	}
	return s.repository.UpdateStatus(id, &old.LifecycleInfo, status)
}

func (s *service) UpdateStatusTx(tx database.Transaction, id uuid.UUID, status LifecycleStatus) error {
	old, err := s.beforeUpdateStatus(id, status)
	if err != nil {
		return err
	}
	return s.repository.UpdateStatusTx(tx, id, &old.LifecycleInfo, status)
}

func (s *service) beforeCreate(lifecycleInfo *LifecycleInfo) error {
	lifecycleInfo.Status = LifecycleStatusCreated
	var parent LifecycleInfo
	var err error
	if lifecycleInfo.HasParentID() {
		parent, _, err = s.repository.GetInfoById(lifecycleInfo.ParentID)
		if err != nil {
			return err
		}
	}
	if err := s.validator.ValidateCreation(lifecycleInfo, &parent); err != nil {
		return err
	}
	return nil
}

func (s *service) beforeUpdate(lifecycleInfo *LifecycleInfo) error {
	old, found, err := s.repository.GetById(lifecycleInfo.ID)
	if err != nil {
		return err
	}
	if !found {
		return ErrLifecycleNotFound
	}
	proposed := Lifecycle{
		LifecycleInfo: *lifecycleInfo,
		Parent:        old.Parent,
		Children:      old.Children,
	}
	if err := s.validator.ValidateUpdate(&old.LifecycleInfo, &proposed); err != nil {
		return err
	}
	return nil
}

func (s *service) beforeDelete(uuid uuid.UUID) error {
	old, found, err := s.repository.GetById(uuid)
	if err != nil {
		return err
	}
	if !found {
		return ErrLifecycleNotFound
	}
	if err = s.validator.ValidateDeletion(&old); err != nil {
		return err
	}
	return nil
}

func (s *service) beforeUpdateStatus(id uuid.UUID, status LifecycleStatus) (*Lifecycle, error) {
	old, found, err := s.repository.GetById(id)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrLifecycleNotFound
	}
	if err := s.validator.ValidateStatusUpdate(&old, status); err != nil {
		return nil, err
	}
	return &old, nil
}
