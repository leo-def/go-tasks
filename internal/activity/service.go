package activity

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type Service interface {
    Create(activity *Activity, createdBy uuid.UUID) error
    Delete(id uuid.UUID) (bool, error)
    Get(p httpx.PaginationParams) ([]Activity, int64, error)
    GetById(id uuid.UUID) (*Activity, bool, error)
    Update(activity *Activity) error
    UpdatedStatus(id uuid.UUID, status lifecycle.LifecycleStatus) error
}

type service struct {
	repository       Repository
	lifecycleService lifecycle.Service
}

func NewService(repository Repository, lifecycleService lifecycle.Service) Service {
	return &service{repository, lifecycleService}
}

func (s *service) Get(p httpx.PaginationParams) ([]Activity, int64, error) {
	return s.repository.Get(p)
}

func (s *service) GetById(id uuid.UUID) (*Activity, bool, error) {
	return s.repository.GetById(id)
}

func (s *service) Delete(id uuid.UUID) (bool, error) {
	activity, found, err := s.repository.GetActivityInfoById(id)
	if !found || err != nil {
		return found, err
	}
	tx := s.repository.GetConnection().BeginTransaction()
	if err = s.lifecycleService.DeleteTx(tx, activity.LifecycleID); err != nil {
		tx.Rollback()
		return false, err
	}
	ok, err := s.repository.DeleteTx(tx, id)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return false, err
	}
	return ok, nil
}

func (s *service) Update(activity *Activity) error {
	existing, found, err := s.repository.GetById(activity.ID)
	activity.LifecycleID = existing.LifecycleID
	if err != nil {
		return err
	}
	if !found {
		return ErrActivityNotFound
	}
	tx := s.repository.GetConnection().BeginTransaction()
	if !activity.Lifecycle.InitDate.IsZero() || !activity.Lifecycle.DueDate.IsZero() {
		lifecycleInfo := lifecycle.LifecycleInfo{
			ID:       existing.LifecycleID,
			InitDate: activity.Lifecycle.InitDate,
			DueDate:  activity.Lifecycle.DueDate,
			Status:   existing.Lifecycle.Status,
		}
		if err := s.lifecycleService.UpdateTx(tx, &lifecycleInfo); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := s.repository.UpdateInfoTx(tx, &activity.ActivityInfo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	if updated, found, err := s.repository.GetById(activity.ID); err != nil {
		return err
	} else if found {
		*activity = *updated
	}
	return nil
}

func (s *service) Create(activity *Activity, createdBy uuid.UUID) error {
	activity.CreatedByID = createdBy
	tx := s.repository.GetConnection().BeginTransaction()
	lifecycle := lifecycle.LifecycleInfo{
		InitDate: activity.Lifecycle.InitDate,
		DueDate:  activity.Lifecycle.DueDate,
	}
	if err := s.lifecycleService.CreateTx(tx, &lifecycle); err != nil {
		tx.Rollback()
		return err
	}
	activity.LifecycleID = lifecycle.ID
	if err := s.repository.CreateInfoTx(tx, &activity.ActivityInfo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	if created, found, err := s.repository.GetById(activity.ID); err != nil {
		return err
	} else if found {
		*activity = *created
	}
	return nil
}

func (s *service) UpdatedStatus(id uuid.UUID, status lifecycle.LifecycleStatus) error {
	activity, found, err := s.repository.GetById(id)
	if err != nil {
		return err
	}
	if !found {
		return ErrActivityNotFound
	}
	tx := s.repository.GetConnection().BeginTransaction()
	if err := s.lifecycleService.UpdateStatusTx(tx, activity.LifecycleID, status); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
