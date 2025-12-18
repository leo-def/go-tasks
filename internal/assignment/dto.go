package assignment

import (
	"go-tasks/internal/collaborator"
	"time"

	"github.com/google/uuid"
)

type AssignmentInfoDTO struct {
	ID            uuid.UUID `json:"id" example:"523e4567-e89b-12d3-a456-426614174000"`
	AssignedToID  uuid.UUID `json:"assigned_to_id" example:"623e4567-e89b-12d3-a456-426614174000"`
	AssignerID    uuid.UUID `json:"assigner_id" example:"723e4567-e89b-12d3-a456-426614174000"`
	TaskID        uuid.UUID `json:"task_id" example:"823e4567-e89b-12d3-a456-426614174000"`
	AssignDate    time.Time `json:"assign_date" example:"2025-01-01T10:00:00Z"`
	AssignEndDate time.Time `json:"assign_end_date" example:"2025-01-15T18:00:00Z"`
	Active        bool      `json:"active" example:"true"`
}

type AssignmentDTO struct {
	AssignmentInfoDTO
	AssignedTo AssignmentParticipationDTO
	Assigner   AssignmentCollaboratorDTO
	Task       AssignmentTaskDTO
}

type TaskAssignmentDTO struct {
	AssignmentInfoDTO
	AssignedTo AssignmentParticipationDTO `json:"assigned_to"`
	Assigner   AssignmentCollaboratorDTO  `json:"assigner"`
}

type AssignmentParticipationDTO struct {
	ID             uuid.UUID `json:"id" example:"923e4567-e89b-12d3-a456-426614174000"`
	CollaboratorID uuid.UUID `json:"collaborator_id" example:"a23e4567-e89b-12d3-a456-426614174000"`
	ActivityID     uuid.UUID `json:"activity_id" example:"b23e4567-e89b-12d3-a456-426614174000"`

	Collaborator AssignmentCollaboratorDTO `json:"collaborator"`
}

type AssignmentCollaboratorDTO struct {
	ID        uuid.UUID                     `json:"id" example:"c23e4567-e89b-12d3-a456-426614174000"`
	CompanyID uuid.UUID                     `json:"company_id" example:"d23e4567-e89b-12d3-a456-426614174000"`
	AccountID uuid.UUID                     `json:"account_id" example:"e23e4567-e89b-12d3-a456-426614174000"`
	Role      collaborator.CollaboratorRole `json:"role" example:"MANAGER"`
	Active    bool                          `json:"active" example:"true"`

	Account AssignmentAccountDTO `json:"account"`
}

type AssignmentTaskDTO struct {
	ID          uuid.UUID `json:"id" example:"f23e4567-e89b-12d3-a456-426614174000"`
	Title       string    `json:"title" example:"Prepare invoice"`
	ActivityID  uuid.UUID `json:"activity_id" example:"013e4567-e89b-12d3-a456-426614174000"`
	LifecycleID uuid.UUID `json:"lifecycle_id" example:"113e4567-e89b-12d3-a456-426614174000"`
}

type AssignmentAccountDTO struct {
	ID       uuid.UUID `json:"id" example:"213e4567-e89b-12d3-a456-426614174000"`
	Username string    `json:"username" example:"asmith"`
	Name     string    `json:"name" example:"Alice Smith"`
	Phone    string    `json:"phone" example:"+1-202-555-0175"`
	Email    string    `json:"email" example:"alice.smith@example.com"`
}

type PaginatedTaskAssignmentDTO struct {
	Items []TaskAssignmentDTO `json:"items"`
	Count int64               `json:"count"`
}

type AssignmentCreateDTO struct {
	AssignedToID  uuid.UUID `json:"assigned_to_id" example:"623e4567-e89b-12d3-a456-426614174000"`
	AssignDate    time.Time `json:"assign_date" example:"2024-12-20T09:00:00Z"`
	AssignEndDate time.Time `json:"assign_end_date" example:"2025-01-15T17:00:00Z"`
}

type AssignmentParticipationCreateDTO struct {
	AssignDate    time.Time `json:"assign_date" example:"2024-12-20T09:00:00Z"`
	AssignEndDate time.Time `json:"assign_end_date" example:"2025-01-15T17:00:00Z"`
}

type AssignmentCollaboratorCreateDTO struct {
	ActivityID    uuid.UUID `json:"activity_id" example:"b23e4567-e89b-12d3-a456-426614174000"`
	AssignDate    time.Time `json:"assign_date" example:"2024-12-20T09:00:00Z"`
	AssignEndDate time.Time `json:"assign_end_date" example:"2025-01-15T17:00:00Z"`
}

func FromAssignmentCreateDTO(dto AssignmentCreateDTO) Assignment {
	return Assignment{
		AssignmentInfo: AssignmentInfo{
			AssignedToID:  dto.AssignedToID,
			AssignDate:    dto.AssignDate,
			AssignEndDate: dto.AssignEndDate,
		},
	}
}

func FromAssignmentParticipationCreateDTO(dto AssignmentParticipationCreateDTO) Assignment {
	return Assignment{
		AssignmentInfo: AssignmentInfo{
			AssignDate:    dto.AssignDate,
			AssignEndDate: dto.AssignEndDate,
		},
	}
}

func FromAssignmentCollaboratorCreateDTO(dto AssignmentCollaboratorCreateDTO) Assignment {
	return Assignment{
		AssignmentInfo: AssignmentInfo{
			AssignDate:    dto.AssignDate,
			AssignEndDate: dto.AssignEndDate,
		},
	}
}

func ToAssignmentDTO(assignment Assignment) AssignmentDTO {
	return AssignmentDTO{
		AssignmentInfoDTO: ToAssignmentInfoDTO(assignment.AssignmentInfo),
		AssignedTo:        ToAssignmentParticipationDTO(assignment.AssignedTo),
		Assigner:          ToAssignmentCollaboratorDTO(assignment.Assigner),
		Task:              ToAssignmentTaskDTO(assignment.Task),
	}
}

func ToAssignmentDTOs(assignments []Assignment) []AssignmentDTO {
	dtos := make([]AssignmentDTO, 0, len(assignments))
	for _, assignment := range assignments {
		dtos = append(dtos, ToAssignmentDTO(assignment))
	}
	return dtos
}

func ToTaskAssignmentDTO(assignment TaskAssignment) TaskAssignmentDTO {
	return TaskAssignmentDTO{
		AssignmentInfoDTO: ToAssignmentInfoDTO(assignment.AssignmentInfo),
		AssignedTo:        ToAssignmentParticipationDTO(assignment.AssignedTo),
		Assigner:          ToAssignmentCollaboratorDTO(assignment.Assigner),
	}
}

func ToTaskAssignmentDTOs(assignments []TaskAssignment) []TaskAssignmentDTO {
	dtos := make([]TaskAssignmentDTO, 0, len(assignments))
	for _, assignment := range assignments {
		dtos = append(dtos, ToTaskAssignmentDTO(assignment))
	}
	return dtos
}

func ToAssignmentInfoDTO(assignment AssignmentInfo) AssignmentInfoDTO {
	return AssignmentInfoDTO{
		ID:            assignment.ID,
		AssignedToID:  assignment.AssignedToID,
		AssignerID:    assignment.AssignerID,
		TaskID:        assignment.TaskID,
		AssignDate:    assignment.AssignDate,
		AssignEndDate: assignment.AssignEndDate,
		Active:        assignment.Active,
	}
}

func ToAssignmentParticipationDTO(participation AssignmentParticipation) AssignmentParticipationDTO {
	return AssignmentParticipationDTO{
		ID:             participation.ID,
		CollaboratorID: participation.CollaboratorID,
		ActivityID:     participation.ActivityID,
		Collaborator:   ToAssignmentCollaboratorDTO(participation.Collaborator),
	}
}

func ToAssignmentCollaboratorDTO(collaborator AssignmentCollaborator) AssignmentCollaboratorDTO {
	return AssignmentCollaboratorDTO{
		ID:        collaborator.ID,
		CompanyID: collaborator.CompanyID,
		AccountID: collaborator.AccountID,
		Role:      collaborator.Role,
		Active:    collaborator.Active,
		Account:   ToAssignmentAccountDTO(collaborator.Account),
	}
}

func ToAssignmentAccountDTO(account AssignmentAccount) AssignmentAccountDTO {
	return AssignmentAccountDTO{
		ID:       account.ID,
		Username: account.Username,
		Name:     account.Name,
		Phone:    account.Phone,
		Email:    account.Email,
	}
}

func ToAssignmentTaskDTO(task AssignmentTask) AssignmentTaskDTO {
	return AssignmentTaskDTO(task)
}
