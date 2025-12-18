package task

import (
	"go-tasks/internal/lifecycle"
	"time"

	"github.com/google/uuid"
)

// Response

type TaskInfoDTO struct {
	ID          uuid.UUID `json:"id" example:"313e4567-e89b-12d3-a456-426614174100"`
	Title       string    `json:"title" example:"Prepare monthly report"`
	ActivityID  uuid.UUID `json:"activity_id" example:"313e4567-e89b-12d3-a456-426614174101"`
	LifecycleID uuid.UUID `json:"lifecycle_id" example:"313e4567-e89b-12d3-a456-426614174102"`
}

type TaskDTO struct {
	TaskInfoDTO `gorm:"embedded"`
	Lifecycle   TaskLifecycleDTO `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
	Activity    TaskActivityDTO  `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}

type ActivityTaskDTO struct {
	TaskInfo  `gorm:"embedded"`
	Lifecycle TaskLifecycleDTO `json:"lifecycle" gorm:"foreignKey:LifecycleID;reference:ID"`
}

type TaskActivityDTO struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey" example:"413e4567-e89b-12d3-a456-426614174100"`
	Title       string    `json:"title" gorm:"not null" example:"Activity Alpha"`
	LifecycleID uuid.UUID `json:"lifecycle_id" gorm:"not null" example:"413e4567-e89b-12d3-a456-426614174101"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"not null" example:"413e4567-e89b-12d3-a456-426614174102"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"not null" example:"413e4567-e89b-12d3-a456-426614174103"`
	CreatedByID uuid.UUID `json:"created_by_id" gorm:"not null" example:"413e4567-e89b-12d3-a456-426614174104"`
	Active      bool      `json:"active" gorm:"not null" example:"true"`

	Lifecycle TaskLifecycleDTO    `json:"lifecycle" gorm:"foreignKey:LifecycleID;reference:ID"`
	Company   TaskCompanyDTO      `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
	Owner     TaskCollaboratorDTO `json:"owner" gorm:"foreignKey:OwnerID;references:ID"`
	CreatedBy TaskCollaboratorDTO `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
}

type TaskCompanyDTO struct {
	ID    uuid.UUID `json:"id" gorm:"primaryKey" example:"513e4567-e89b-12d3-a456-426614174100"`
	Title string    `json:"title" gorm:"not null" example:"Acme Corp"`
}

type TaskLifecycleDTO struct {
	ID       uuid.UUID             `json:"id" gorm:"primaryKey" example:"613e4567-e89b-12d3-a456-426614174100"`
	InitDate time.Time             `json:"init_date" gorm:"not null" example:"2025-01-01T08:00:00Z"`
	DueDate  time.Time             `json:"due_date" gorm:"not null" example:"2025-01-10T18:00:00Z"`
	Status   string                `json:"status" gorm:"not null" example:"OPEN"`
	Updates  []TaskStatusUpdateDTO `json:"updates" gorm:"foreignKey:LifecycleID;references:ID"`
}

type TaskStatusUpdateDTO struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey" example:"713e4567-e89b-12d3-a456-426614174100"`
	LifecycleID  uuid.UUID `json:"lifecycle_id" gorm:"not null" example:"613e4567-e89b-12d3-a456-426614174100"`
	StatusBefore string    `json:"status_before" gorm:"not null" example:"OPEN"`
	StatusAfter  string    `json:"status_after" gorm:"not null" example:"IN_PROGRESS"`
	UpdateDate   time.Time `json:"update_date" gorm:"not null" example:"2025-01-05T12:00:00Z"`
}

type TaskCollaboratorDTO struct {
    ID        uuid.UUID      `json:"id" gorm:"primaryKey" example:"813e4567-e89b-12d3-a456-426614174100"`
    Role      string         `json:"role" gorm:"not null" example:"MANAGER"`
    AccountID uuid.UUID      `json:"account_id" gorm:"not null" example:"913e4567-e89b-12d3-a456-426614174100"`
    CompanyID uuid.UUID      `json:"company_id" gorm:"not null" example:"a13e4567-e89b-12d3-a456-426614174100"`
    Account   TaskAccountDTO `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

type TaskAccountDTO struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey" example:"b13e4567-e89b-12d3-a456-426614174100"`
	Name     string    `json:"name" gorm:"not null" example:"Jane Doe"`
	Username string    `json:"username" gorm:"not null" example:"jdoe"`
	Email    string    `json:"email" gorm:"not null" example:"jdoe@example.com"`
	Phone    string    `json:"phone" gorm:"not null" example:"+1-202-555-0188"`
	Role     string    `json:"role" gorm:"not null" example:"ADMIN"`
}

type PaginatedActivityTaskDTO struct {
	Items []ActivityTaskDTO `json:"items"`
	Count int64             `json:"count"`
}

// Create Param
type TaskCreateDTO struct {
	Title     string                 `json:"title" example:"Prepare monthly report"`
	Lifecycle TaskLifecycleCreateDTO `json:"lifecycle"`
}

type TaskLifecycleCreateDTO struct {
	InitDate time.Time `json:"init_date" example:"2025-01-01T08:00:00Z"`
	DueDate  time.Time `json:"due_date" example:"2025-01-10T18:00:00Z"`
}

// Response Param
type TaskUpdateDTO struct {
	Title     string                 `json:"title" example:"Prepare monthly report (updated)"`
	Lifecycle TaskLifecycleUpdateDTO `json:"lifecycle"`
}

type TaskLifecycleUpdateDTO struct {
	InitDate time.Time `json:"init_date" example:"2025-01-01T08:00:00Z"`
	DueDate  time.Time `json:"due_date" example:"2025-01-10T18:00:00Z"`
}

// Update Status
type TaskUpdateStatusDTO struct {
	Status lifecycle.LifecycleStatusDTO `json:"status" example:"IN_PROGRESS"`
}

func FromCreateDTO(dto TaskCreateDTO) Task {
	return Task{
		TaskInfo: TaskInfo{
			Title: dto.Title,
		},
		Lifecycle: TaskLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func FromUpdateDTOWithId(dto TaskUpdateDTO, id uuid.UUID) Task {
	return Task{
		TaskInfo: TaskInfo{
			ID:    id,
			Title: dto.Title,
		},
		Lifecycle: TaskLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func ToActivityTaskDTO(task *ActivityTask) ActivityTaskDTO {
	return ActivityTaskDTO{
		TaskInfo: TaskInfo{
			ID:          task.ID,
			Title:       task.Title,
			ActivityID:  task.ActivityID,
			LifecycleID: task.LifecycleID,
		},
		Lifecycle: ToTaskLifecycleDTO(task.Lifecycle),
	}
}

func ToActivityTaskDTOs(tasks []ActivityTask) []ActivityTaskDTO {
	dtos := make([]ActivityTaskDTO, 0, len(tasks))
	for _, task := range tasks {
		dtos = append(dtos, ToActivityTaskDTO(&task))
	}
	return dtos
}

func ToTaskDTO(task *Task) TaskDTO {
	return TaskDTO{
		TaskInfoDTO: TaskInfoDTO{
			ID:          task.ID,
			Title:       task.Title,
			ActivityID:  task.ActivityID,
			LifecycleID: task.LifecycleID,
		},
		Lifecycle: ToTaskLifecycleDTO(task.Lifecycle),
		Activity:  ToTaskActivityDTO(task.Activity),
	}
}

func ToTaskActivityDTO(activity TaskActivity) TaskActivityDTO {
	return TaskActivityDTO{
		ID:          activity.ID,
		Title:       activity.Title,
		LifecycleID: activity.LifecycleID,
		CompanyID:   activity.CompanyID,
		OwnerID:     activity.OwnerID,
		CreatedByID: activity.CreatedByID,
		Active:      activity.Active,
		Lifecycle:   ToTaskLifecycleDTO(activity.Lifecycle),
		Company:     ToTaskCompanyDTO(activity.Company),
		Owner:       ToTaskCollaboratorDTO(activity.Owner),
		CreatedBy:   ToTaskCollaboratorDTO(activity.CreatedBy),
	}
}

func ToTaskLifecycleDTO(lifecycle TaskLifecycle) TaskLifecycleDTO {
	return TaskLifecycleDTO{
		ID:       lifecycle.ID,
		InitDate: lifecycle.InitDate,
		DueDate:  lifecycle.DueDate,
		Status:   lifecycle.Status,
		Updates:  ToTaskStatusUpdateDTOs(lifecycle.Updates),
	}
}

func ToTaskStatusUpdateDTO(statusUpdate TaskStatusUpdate) TaskStatusUpdateDTO {
	return TaskStatusUpdateDTO(statusUpdate)
}

func ToTaskStatusUpdateDTOs(status []TaskStatusUpdate) []TaskStatusUpdateDTO {
	dtos := make([]TaskStatusUpdateDTO, 0, len(status))
	for _, item := range status {
		dtos = append(dtos, ToTaskStatusUpdateDTO(item))
	}
	return dtos
}

func ToTaskCompanyDTO(company TaskCompany) TaskCompanyDTO {
	return TaskCompanyDTO(company)
}

func ToTaskCollaboratorDTO(collaborator TaskCollaborator) TaskCollaboratorDTO {
    return TaskCollaboratorDTO{
        ID:        collaborator.ID,
        Role:      collaborator.Role,
        AccountID: collaborator.AccountID,
        CompanyID: collaborator.CompanyID,
        Account:   ToTaskAccountDTO(collaborator.Account),
    }
}

func ToTaskAccountDTO(account TaskAccount) TaskAccountDTO {
	return TaskAccountDTO(account)
}
