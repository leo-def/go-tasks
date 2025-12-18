package assignment

import (
	"go-tasks/internal/collaborator"
	"time"

	"github.com/google/uuid"
)

/*
Assignment -> AssignedTo
Assignment -> AssignedTo -> Collaborator
Assignment -> AssignedTo -> Collaborator -> Account
Assignment -> Assigner
Assignment -> Assigner -> Account
Assignment -> Task
*/
type AssignmentInfo struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AssignedToID  uuid.UUID `json:"assigned_to_id" gorm:"type:uuid"`
	AssignerID    uuid.UUID `json:"assigner_id" gorm:"type:uuid"`
	TaskID        uuid.UUID `json:"task_id" gorm:"type:uuid"`
	Active        bool      `json:"active" gorm:"not null;default:true"`
	AssignDate    time.Time `json:"assign_date"`
	AssignEndDate time.Time `json:"assign_end_date"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     time.Time `json:"deleted_at" gorm:"index"`
}

func (AssignmentInfo) TableName() string { return "assignments" }

type Assignment struct {
	AssignmentInfo `gorm:"embedded"`
	AssignedTo     AssignmentParticipation `json:"assigned_to" gorm:"foreignKey:AssignedToID;references:ID"`
	Assigner       AssignmentCollaborator  `json:"assigner" gorm:"foreignKey:AssignerID;references:ID"`
	Task           AssignmentTask          `json:"task" gorm:"foreignKey:TaskID;references:ID"`
}

type TaskAssignment struct {
	AssignmentInfo `gorm:"embedded"`
	AssignedTo     AssignmentParticipation `json:"assigned_to" gorm:"foreignKey:AssignedToID;references:ID"`
	Assigner       AssignmentCollaborator  `json:"assigner" gorm:"foreignKey:AssignerID;references:ID"`
}

func (TaskAssignment) TableName() string { return "assignments" }

type AssignmentParticipation struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CollaboratorID uuid.UUID `json:"collaborator_id" gorm:"type:uuid;not null"`
	ActivityID     uuid.UUID `json:"activity_id" gorm:"type:uuid;not null"`

	Collaborator AssignmentCollaborator `json:"collaborator" gorm:"foreignKey:CollaboratorID;references:ID"`
}

func (AssignmentParticipation) TableName() string { return "participations" }

type AssignmentCollaborator struct {
	ID        uuid.UUID                     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID                     `json:"company_id" gorm:"type:uuid;not null"`
	AccountID uuid.UUID                     `json:"account_id" gorm:"type:uuid;not null"`
	Role      collaborator.CollaboratorRole `json:"role" gorm:"not null"`
	Active    bool                          `json:"active" gorm:"not null;default:true"`

	Account AssignmentAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (AssignmentCollaborator) TableName() string { return "collaborators" }

type AssignmentTask struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	ActivityID  uuid.UUID `json:"activity_id" gorm:"type:uuid;not null"`
	LifecycleID uuid.UUID `json:"lifecycle_id" gorm:"type:uuid;not null"`
}

func (AssignmentTask) TableName() string { return "tasks" }

type AssignmentAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"not null"`
	Name     string    `json:"name" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Role     string    `json:"role" gorm:"not null"`
}

func (AssignmentAccount) TableName() string { return "accounts" }
