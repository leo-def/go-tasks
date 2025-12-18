package task

import (
	"go-tasks/internal/activity"
	"go-tasks/internal/assignment"
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceActivity interface {
    Create(task *Task, activityID uuid.UUID) error
    Delete(id uuid.UUID, activityID uuid.UUID) (bool, error)
    Get(p httpx.PaginationParams, activityID uuid.UUID) ([]ActivityTask, int64, error)
    GetById(id uuid.UUID, activityID uuid.UUID) (*ActivityTask, bool, error)
    GetByStatus(status []lifecycle.LifecycleStatus, activityID uuid.UUID) ([]ActivityTask, int64, error)
    Perform(id uuid.UUID, collaboratorId uuid.UUID, status lifecycle.LifecycleStatus, activityID uuid.UUID) error
    Update(task *Task, activityID uuid.UUID) error
    UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, activityID uuid.UUID) error
}

type serviceActivity struct {
    repository        Repository
    lifecycleService  lifecycle.Service
    assignmentService assignment.ServiceTask
    activityService   activity.Service
}

func NewService(repository Repository, lifecycleService lifecycle.Service, assignmentService assignment.ServiceTask, activityService activity.Service) ServiceActivity {
	return &serviceActivity{repository: repository, lifecycleService: lifecycleService, assignmentService: assignmentService, activityService: activityService}
}

func (s *serviceActivity) Create(task *Task, activityID uuid.UUID) error {
    task.ActivityID = activityID
    activity, found, err := s.activityService.GetById(task.ActivityID)
    if err != nil {
        return err
    } else if !found {
        return ErrActivityNotFoundForTask
    }

	tx := s.repository.GetConnection().BeginTransaction()
	lifecycle := lifecycle.LifecycleInfo{
		ParentID: activity.LifecycleID,
		InitDate: task.Lifecycle.InitDate,
		DueDate:  task.Lifecycle.DueDate,
	}
	if err := s.lifecycleService.CreateTx(tx, &lifecycle); err != nil {
		tx.Rollback()
		return err
	}
	task.LifecycleID = lifecycle.ID
	if err := s.repository.CreateInfoTx(tx, &task.TaskInfo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	if created, found, err := s.repository.GetById(task.ID); err != nil {
		return err
	} else if found {
		*task = *created
    }
    return nil
}

func (s *serviceActivity) Delete(id uuid.UUID, activityID uuid.UUID) (bool, error) {
    task, found, err := s.repository.GetActivityTaskById(id, activityID)
    if err != nil {
        return false, err
    }
    if !found {
        return false, nil
    }
    tx := s.repository.GetConnection().BeginTransaction()
    if err = s.lifecycleService.DeleteTx(tx, task.LifecycleID); err != nil {
        tx.Rollback()
        return false, err
    }
    if ok, err := s.repository.DeleteTx(tx, id); err != nil {
        tx.Rollback()
        return false, err
    } else if err := tx.Commit(); err != nil {
        tx.Rollback()
        return false, err
    } else {
        return ok, nil
    }
}

func (s *serviceActivity) Get(p httpx.PaginationParams, activityID uuid.UUID) ([]ActivityTask, int64, error) {
    return s.repository.GetByActivityId(activityID, p)
}

func (s *serviceActivity) GetById(id uuid.UUID, activityID uuid.UUID) (*ActivityTask, bool, error) {
    if activityTask, found, err := s.repository.GetActivityTaskById(id, activityID); !found || err != nil {
        return nil, found, err
    } else {
        return activityTask, true, nil
    }
}

func (s *serviceActivity) GetByStatus(status []lifecycle.LifecycleStatus, activityID uuid.UUID) ([]ActivityTask, int64, error) {
    return s.repository.GetByActivityIdInStatus(activityID, status)
}

func (s *serviceActivity) Perform(id uuid.UUID, collaboratorId uuid.UUID, status lifecycle.LifecycleStatus, activityID uuid.UUID) error {
    if _, found, err := s.assignmentService.GetCollaboratorTask(collaboratorId, id); err != nil {
        return err
    } else if !found {
        return ErrAssignmentNotFound
    }

	task, found, err := s.repository.GetById(id)
	if err != nil {
		return err
	} else if !found {
		return ErrTaskNotFound
	}
	tx := s.repository.GetConnection().BeginTransaction()
	if err := s.lifecycleService.UpdateStatusTx(tx, task.LifecycleID, status); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
    }
    return nil
}

func (s *serviceActivity) Update(task *Task, activityID uuid.UUID) error {
    existing, found, err := s.repository.GetById(task.ID)
    if err != nil {
        return err
    } else if !found {
        return ErrTaskNotFound
    } else if existing.ActivityID != activityID {
        return ErrActivityNotFoundForTask
    }
    tx := s.repository.GetConnection().BeginTransaction()
    if !task.Lifecycle.InitDate.IsZero() || !task.Lifecycle.DueDate.IsZero() {
        lifecycleInfo := lifecycle.LifecycleInfo{
            ID:       existing.LifecycleID,
            ParentID: existing.Activity.LifecycleID,
            InitDate: task.Lifecycle.InitDate,
            DueDate:  task.Lifecycle.DueDate,
            Status:   lifecycle.LifecycleStatus(existing.Lifecycle.Status),
        }
        if err := s.lifecycleService.UpdateTx(tx, &lifecycleInfo); err != nil {
            tx.Rollback()
            return err
        }
    }
    if err := s.repository.UpdateInfoTx(tx, &task.TaskInfo); err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }
    if updated, found, err := s.repository.GetById(task.ID); err != nil {
        return err
    } else if found {
        *task = *updated
    }
    return nil
}

func (s *serviceActivity) UpdateStatus(id uuid.UUID, status lifecycle.LifecycleStatus, activityID uuid.UUID) error {
    task, found, err := s.repository.GetById(id)
    if err != nil {
        return err
    } else if !found {
        return ErrTaskNotFound
    }
    tx := s.repository.GetConnection().BeginTransaction()
    if err := s.lifecycleService.UpdateStatusTx(tx, task.LifecycleID, status); err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }
    return nil
}
