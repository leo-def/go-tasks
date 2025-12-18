package lifecycle

import (
	"time"

	"github.com/google/uuid"
)

type LifecycleStatus string

const (
	LifecycleStatusCreated    LifecycleStatus = "CREATED"
	LifecycleStatusInProgress LifecycleStatus = "IN_PROGRESS"
	LifecycleStatusBlocked    LifecycleStatus = "BLOCKED"
	LifecycleStatusCompleted  LifecycleStatus = "COMPLETED"
	LifecycleStatusCancelled  LifecycleStatus = "CANCELLED"
)

type LifecycleInfo struct {
	ID       uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InitDate time.Time       `json:"init_date" gorm:"not null"`
	DueDate  time.Time       `json:"due_date" gorm:"not null"`
    ParentID uuid.UUID       `json:"parent_id" gorm:"type:uuid"`
	Status   LifecycleStatus `json:"status" gorm:"not null"`
}

func (LifecycleInfo) TableName() string {
	return "lifecycles"
}

func (lifecycle *LifecycleInfo) HasID() bool {
	return lifecycle.ID != uuid.Nil
}

func (lifecycle *LifecycleInfo) HasParentID() bool {
	return lifecycle.ParentID != uuid.Nil
}

type Lifecycle struct {
	LifecycleInfo `gorm:"embedded"`
	Updates       []LifecycleStatusUpdate `json:"updates" gorm:"foreignKey:LifecycleID;references:ID"`
	Parent        LifecycleInfo           `json:"parent" gorm:"foreignKey:ParentID;references:ID"`
	Children      []LifecycleInfo         `json:"children" gorm:"foreignKey:ParentID;references:ID"`
}

func (Lifecycle) TableName() string {
	return "lifecycles"
}

func (lifecycle *Lifecycle) HasID() bool {
	return lifecycle.LifecycleInfo.HasID()
}

func (lifecycle *Lifecycle) HasParentID() bool {
	return lifecycle.LifecycleInfo.HasParentID()
}

type LifecycleStatusUpdate struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LifecycleID  uuid.UUID       `json:"lifecycle_id" gorm:"not null"`
	StatusBefore LifecycleStatus `json:"status_before" gorm:"not null"`
	StatusAfter  LifecycleStatus `json:"status_after" gorm:"not null"`
	UpdateDate   time.Time       `json:"update_date" gorm:"not null"`
}

func (LifecycleStatusUpdate) TableName() string {
	return "status_updates"
}
