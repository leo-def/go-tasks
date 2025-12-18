package rating

import (
	"go-tasks/internal/collaborator"
	"time"

	"github.com/google/uuid"
)

type RatingInfo struct {
    ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Rate            float64   `json:"rate" gorm:"not null"`
    Comment         string    `json:"comment" gorm:"type:text"`
    RateType        string    `json:"rate_type" gorm:"type:text"`
    ParticipationID uuid.UUID `json:"participation_id" gorm:"type:uuid;not null"`
    CollaboratorID  uuid.UUID `json:"collaborator_id" gorm:"type:uuid;not null"`
    CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt       time.Time `json:"deleted_at" gorm:"index"`
}

func (RatingInfo) TableName() string {
	return "ratings"
}

type Rating struct {
	RatingInfo    `gorm:"embedded"`
	Participation RatingParticipation `json:"participation" gorm:"foreignKey:ParticipationID;references:ID"`
	Collaborator  RatingCollaborator  `json:"collaborator" gorm:"foreignKey:CollaboratorID;references:ID"`
}

type RatingParticipation struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CollaboratorID uuid.UUID `json:"collaborator_id" gorm:"type:uuid;not null"`
	ActivityID     uuid.UUID `json:"activity_id" gorm:"type:uuid;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"index"`
}

func (RatingParticipation) TableName() string {
	return "participations"
}

type RatingCollaborator struct {
	ID        uuid.UUID                     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID                     `json:"company_id" gorm:"type:uuid;not null"`
	AccountID uuid.UUID                     `json:"account_id" gorm:"type:uuid;not null"`
	Role      collaborator.CollaboratorRole `json:"role" gorm:"not null"`
	Active    bool                          `json:"active" gorm:"not null;default:true"`
	CreatedAt time.Time                     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time                     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time                     `json:"deleted_at" gorm:"index"`
}

func (RatingCollaborator) TableName() string {
	return "collaborators"
}
