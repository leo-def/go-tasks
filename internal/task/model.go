package task

import (
	"time"

	"github.com/google/uuid"
)

type TaskInfo struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	ActivityID  uuid.UUID `json:"activity_id" gorm:"type:uuid;not null"`
	LifecycleID uuid.UUID `json:"lifecycle_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`
}

func (TaskInfo) TableName() string {
	return "tasks"
}

type Task struct {
	TaskInfo  `gorm:"embedded"`
	Lifecycle TaskLifecycle `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
	Activity  TaskActivity  `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}

type ActivityTask struct {
	TaskInfo  `gorm:"embedded"`
	Lifecycle TaskLifecycle `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
}

func (ActivityTask) TableName() string {
	return "tasks"
}

type TaskActivity struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	LifecycleID uuid.UUID `json:"lifecycle_id" gorm:"not null"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"not null"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"not null"`
	CreatedByID uuid.UUID `json:"created_by_id" gorm:"not null"`
	Active      bool      `json:"active" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"index"`

	Lifecycle TaskLifecycle    `json:"lifecycle" gorm:"foreignKey:LifecycleID;reference:ID"`
	Company   TaskCompany      `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
	Owner     TaskCollaborator `json:"owner" gorm:"foreignKey:OwnerID;references:ID"`
	CreatedBy TaskCollaborator `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
}

func (TaskActivity) TableName() string {
	return "activities"
}

type TaskCompany struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title string    `json:"title" gorm:"not null"`
}

func (TaskCompany) TableName() string {
	return "companies"
}

type TaskLifecycle struct {
	ID       uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InitDate time.Time          `json:"init_date" gorm:"not null"`
	DueDate  time.Time          `json:"due_date" gorm:"not null"`
	Status   string             `json:"status" gorm:"not null"`
	Updates  []TaskStatusUpdate `json:"updates" gorm:"foreignKey:LifecycleID;references:ID"`
}

func (TaskLifecycle) TableName() string {
	return "lifecycles"
}

type TaskStatusUpdate struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LifecycleID  uuid.UUID `json:"lifecycle_id" gorm:"not null"`
	StatusBefore string    `json:"status_before" gorm:"not null"`
	StatusAfter  string    `json:"status_after" gorm:"not null"`
	UpdateDate   time.Time `json:"update_date" gorm:"not null"`
}

func (TaskStatusUpdate) TableName() string {
	return "status_updates"
}

type TaskCollaborator struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Role      string      `json:"role" gorm:"not null"`
	AccountID uuid.UUID   `json:"account_id" gorm:"not null"`
	CompanyID uuid.UUID   `json:"company_id" gorm:"not null"`
	Account   TaskAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (TaskCollaborator) TableName() string {
	return "collaborators"
}

type TaskAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `json:"name" gorm:"not null"`
	Username string    `json:"username" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Role     string    `json:"role" gorm:"not null"`
}

func (TaskAccount) TableName() string {
	return "accounts"
}
