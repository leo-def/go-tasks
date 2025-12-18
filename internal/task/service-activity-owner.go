package task

import (
	"go-tasks/internal/activity"
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceActivityOwner interface {
    Create(task *Task, activityID, ownerID uuid.UUID) error
    Delete(id, activityID, ownerID uuid.UUID) (bool, error)
    Get(p httpx.PaginationParams, activityID, ownerID uuid.UUID) ([]ActivityTask, int64, error)
    GetById(id, activityID, ownerID uuid.UUID) (*ActivityTask, bool, error)
    Update(task *Task, activityID, ownerID uuid.UUID) error
    UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, activityID, ownerID uuid.UUID) error
}

type serviceActivityOwner struct {
	base            ServiceActivity
	activityService activity.Service
}

func NewServiceActivityOwner(base ServiceActivity, activityService activity.Service) ServiceActivityOwner {
	return &serviceActivityOwner{base: base, activityService: activityService}
}

func (s *serviceActivityOwner) Create(task *Task, activityID, ownerID uuid.UUID) error {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return err
    }
    return s.base.Create(task, activityID)
}

func (s *serviceActivityOwner) Delete(id, activityID, ownerID uuid.UUID) (bool, error) {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return false, err
    }
    return s.base.Delete(id, activityID)
}

func (s *serviceActivityOwner) Get(p httpx.PaginationParams, activityID, ownerID uuid.UUID) ([]ActivityTask, int64, error) {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return nil, 0, err
    }
    return s.base.Get(p, activityID)
}

func (s *serviceActivityOwner) GetById(id, activityID, ownerID uuid.UUID) (*ActivityTask, bool, error) {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return nil, false, err
    }
    return s.base.GetById(id, activityID)
}

func (s *serviceActivityOwner) Update(task *Task, activityID, ownerID uuid.UUID) error {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return err
    }
    return s.base.Update(task, activityID)
}

func (s *serviceActivityOwner) UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, activityID, ownerID uuid.UUID) error {
    if err := s.verifyOwnership(activityID, ownerID); err != nil {
        return err
    }
    return s.base.UpdateStatus(id, status, activityID)
}

func (s *serviceActivityOwner) verifyOwnership(activityID, ownerID uuid.UUID) error {
	a, found, err := s.activityService.GetById(activityID)
	if err != nil {
		return err
	}
	if !found {
		return activity.ErrActivityNotFound
	}
	if a.Owner.ID != ownerID {
		return activity.ErrForbiddenNotOwner
	}
	return nil
}
