package activity

import (
	"go-tasks/internal/lifecycle"
	"time"

	"github.com/google/uuid"
)

/*
CompanyActivity:
{
	"title": "",
	"lifecycle": {
		"init_date": "",
		"due_date": "",
		"status": "",
	},
	"owner": {
		"title": "",
		"role": "",
	},
	"created_by": {
		"title": "",
		"role": "",
	},
}
*/

type ActivityInfo struct {
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
}

func (ActivityInfo) TableName() string {
	return "activities"
}

type Activity struct {
	ActivityInfo `gorm:"embedded"`
	Company      ActivityCompany      `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
	Lifecycle    ActivityLifecycle    `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
	Owner        ActivityCollaborator `json:"owner" gorm:"foreignKey:OwnerID;references:ID"`
	CreatedBy    ActivityCollaborator `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
}

type CompanyActivity struct {
	ActivityInfo `gorm:"embedded"`
	Lifecycle    ActivityLifecycle    `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
	Owner        ActivityCollaborator `json:"owner" gorm:"foreignKey:OwnerID;references:ID"`
	CreatedBy    ActivityCollaborator `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
}

type OwnActivity struct {
	ActivityInfo `gorm:"embedded"`
	Lifecycle    ActivityLifecycle    `json:"lifecycle" gorm:"foreignKey:LifecycleID;references:ID"`
	CreatedBy    ActivityCollaborator `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
}

func (CompanyActivity) TableName() string {
	return "activities"
}

type ActivityCompany struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title string    `json:"title" gorm:"not null"`
}

func (ActivityCompany) TableName() string {
	return "companies"
}

type ActivityLifecycle struct {
	ID       uuid.UUID                 `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InitDate time.Time                 `json:"init_date" gorm:"not null"`
	DueDate  time.Time                 `json:"due_date" gorm:"not null"`
	Status   lifecycle.LifecycleStatus `json:"status" gorm:"not null"`
	Updates  []ActivityStatusUpdate    `json:"updates" gorm:"foreignKey:LifecycleID;references:ID"`
}

func (ActivityLifecycle) TableName() string {
	return "lifecycles"
}

type ActivityStatusUpdate struct {
	ID           uuid.UUID                 `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LifecycleID  uuid.UUID                 `json:"lifecycle_id" gorm:"not null"`
	StatusBefore lifecycle.LifecycleStatus `json:"status_before" gorm:"not null"`
	StatusAfter  lifecycle.LifecycleStatus `json:"status_after" gorm:"not null"`
	UpdateDate   time.Time                 `json:"update_date" gorm:"not null"`
}

func (ActivityStatusUpdate) TableName() string {
	return "status_updates"
}

type ActivityCollaborator struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Role      string          `json:"role" gorm:"not null"`
	AccountID uuid.UUID       `json:"account_id" gorm:"not null"`
	CompanyID uuid.UUID       `json:"company_id" gorm:"not null"`
	Account   ActivityAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (ActivityCollaborator) TableName() string {
	return "collaborators"
}

type ActivityAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `json:"name" gorm:"not null"`
	Username string    `json:"username" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Role     string    `json:"role" gorm:"not null"`
}

func (ActivityAccount) TableName() string {
	return "accounts"
}
