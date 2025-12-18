package assignment

import (
	"go-tasks/internal/participation"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
)

type ServiceTask interface {
	GetByTaskId(p httpx.PaginationParams, taskId uuid.UUID) ([]TaskAssignment, int64, error)
	GetCollaboratorTask(collaboratorId, taskId uuid.UUID) (*TaskAssignment, bool, error)
	GetById(id uuid.UUID, taskID uuid.UUID) (*Assignment, bool, error)
	Delete(id uuid.UUID, taskID uuid.UUID) (bool, error)
	Deactivate(id uuid.UUID, taskID uuid.UUID) (bool, error)
	Create(assignment *Assignment, assignerID, taskID uuid.UUID) error
	CreateForParticipation(assignment *Assignment, assignerID, participationID, taskID uuid.UUID) error
	CreateForCollaborator(assignment *Assignment, assignerID, activityID, collaboratorID, taskID uuid.UUID) error
}

type serviceTask struct {
	repository           Repository
	participationService participation.Service
}

func NewServiceTask(r Repository, participationService participation.Service) ServiceTask {
	return &serviceTask{r, participationService}
}

func (s *serviceTask) GetByTaskId(p httpx.PaginationParams, taskId uuid.UUID) ([]TaskAssignment, int64, error) {
	return s.repository.GetByTaskId(taskId, p)
}

func (s *serviceTask) GetCollaboratorTask(collaboratorId, taskId uuid.UUID) (*TaskAssignment, bool, error) {
	return s.repository.GetByTaskAndCollaboratorId(taskId, collaboratorId)
}

func (s *serviceTask) GetById(id uuid.UUID, taskID uuid.UUID) (*Assignment, bool, error) {
	err := s.VerifyAssignerTask(id, taskID)
	if err != nil {
		return nil, false, err
	}
	return s.repository.GetById(id)
}

func (s *serviceTask) Delete(id uuid.UUID, taskID uuid.UUID) (bool, error) {
	err := s.VerifyAssignerTask(id, taskID)
	if err != nil {
		return false, err
	}
	return s.repository.Delete(id)
}

func (s *serviceTask) Deactivate(id uuid.UUID, taskID uuid.UUID) (bool, error) {
	err := s.VerifyAssignerTask(id, taskID)
	if err != nil {
		return false, err
	}
	return s.repository.Deactivate(id)
}

func (s *serviceTask) Create(assignment *Assignment, assignerID, taskID uuid.UUID) error {
	assignment.AssignerID = assignerID
	assignment.TaskID = taskID
	assignment.Active = true // Ensure new assignment is active

	// Use transaction to ensure atomicity
	tx := s.repository.GetConnection().BeginTransaction()

	// First, deactivate any existing active assignments for this task
	_, err := s.repository.DeactivateByTaskIdTx(tx, taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then create the new assignment
	if err := s.repository.CreateTx(tx, assignment); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *serviceTask) CreateForParticipation(assignment *Assignment, assignerID, participationID, taskID uuid.UUID) error {
	assignment.AssignerID = assignerID
	assignment.TaskID = taskID
	assignment.AssignedToID = participationID
	assignment.Active = true

	// Use transaction to ensure atomicity
	tx := s.repository.GetConnection().BeginTransaction()

	// First, deactivate any existing active assignments for this task
	_, err := s.repository.DeactivateByTaskIdTx(tx, taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then create the new assignment
	if err := s.repository.CreateTx(tx, assignment); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *serviceTask) CreateForCollaborator(assignment *Assignment, assignerID, activityID, collaboratorID, taskID uuid.UUID) error {
	tx := s.repository.GetConnection().BeginTransaction()

	// Get or create participation for the collaborator in the activity
	participation, found, err := s.participationService.GetOrCreateByActivityAndCollaboratorTx(tx, activityID, collaboratorID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !found {
		println("participation not found, created new activityID: ", activityID.String(), " collaboratorID: ", collaboratorID.String())
	}

	assignment.AssignerID = assignerID
	assignment.TaskID = taskID
	assignment.AssignedToID = participation.ID
	assignment.Active = true

	// Deactivate any existing active assignments for this task
	_, err = s.repository.DeactivateByTaskIdTx(tx, taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create the new assignment
	if err := s.repository.CreateTx(tx, assignment); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s serviceTask) VerifyAssignerTask(id uuid.UUID, taskID uuid.UUID) error {
	found, err := s.repository.CheckByIdAndTaskId(id, taskID)
	if err != nil {
		return err
	}
	if !found {
		return ErrNotAssignerTask
	}
	return nil
}
