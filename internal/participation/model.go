package participation

import (
	"time"

	"github.com/google/uuid"
)

type ParticipationInfo struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CollaboratorID uuid.UUID `json:"collaborator_id" gorm:"type:uuid;not null"`
	ActivityID     uuid.UUID `json:"activity_id" gorm:"type:uuid;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"index"`
}

func (ParticipationInfo) TableName() string {
	return "participations"
}

type Participation struct {
	ParticipationInfo `gorm:"embedded"`
	Collaborator      ParticipationCollaborator `json:"collaborator" gorm:"foreignKey:CollaboratorID;references:ID"`
	Activity          ParticipationActivity     `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}

type CollaboratorParticipation struct {
	ParticipationInfo `gorm:"embedded"`
	Activity          ParticipationActivity `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}

func (CollaboratorParticipation) TableName() string {
	return "participations"
}

type ActivityParticipation struct {
	ParticipationInfo `gorm:"embedded"`
	Collaborator      ParticipationCollaborator `json:"collaborator" gorm:"foreignKey:CollaboratorID;references:ID"`
}

func (ActivityParticipation) TableName() string {
	return "participations"
}

type ParticipationCollaborator struct {
	ID        uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string               `json:"title" gorm:"not null"`
	AccountID uuid.UUID            `json:"account_id" gorm:"type:uuid;not null"`
	Account   ParticipationAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (ParticipationCollaborator) TableName() string {
	return "collaborators"
}

type ParticipationAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"not null"`
	Name     string    `json:"name" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
}

func (ParticipationAccount) TableName() string {
	return "accounts"
}

type ParticipationActivity struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title string    `json:"title" gorm:"not null"`
}

func (ParticipationActivity) TableName() string { return "activities" }
