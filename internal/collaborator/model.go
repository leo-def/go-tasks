package collaborator

import (
	"time"

	"github.com/google/uuid"
)

type CollaboratorRole string

const (
	CollaboratorRoleOwner   CollaboratorRole = "OWNER"
	CollaboratorRoleManager CollaboratorRole = "MANAGER"
	CollaboratorRoleOps     CollaboratorRole = "OPS"
)

type CollaboratorInfo struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID        `json:"company_id" gorm:"type:uuid;not null"`
	AccountID uuid.UUID        `json:"account_id" gorm:"type:uuid;not null"`
	Role      CollaboratorRole `json:"role" gorm:"not null"`
	Active    bool             `json:"active" gorm:"not null;default:true"`
	CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time        `json:"deleted_at" gorm:"index"`
}

func (CollaboratorInfo) TableName() string { return "collaborators" }

type CollaboratorInfoWithToken struct {
	CollaboratorInfo `gorm:"embedded"`
	ActivationToken  string `json:"activation_token"`
}

func (CollaboratorInfoWithToken) TableName() string { return "collaborators" }

type Collaborator struct {
	CollaboratorInfo `gorm:"embedded"`
	Company          CollaboratorCompany `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
	Account          CollaboratorAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

// CollaboratorCompany mirrors key fields of company.Company to avoid import cycles.
// It maps to the same DB table via TableName.
type CollaboratorCompany struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title string    `json:"title" gorm:"not null"`
}

func (CollaboratorCompany) TableName() string { return "companies" }

// CollaboratorAccount mirrors key fields of account.Account for responses.
// It maps to the same DB table via TableName.
type CollaboratorAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string    `json:"username" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Name     string    `json:"name" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Email    string    `json:"email" gorm:"unique;not null"`
	Role     string    `json:"role" gorm:"not null"`
}

func (CollaboratorAccount) TableName() string { return "accounts" }

type AccountCollaborator struct {
	CollaboratorInfo `gorm:"embedded"`
	Company          CollaboratorCompany `json:"company" gorm:"foreignKey:CompanyID;references:ID"`
}

func (AccountCollaborator) TableName() string { return "collaborators" }

type CompanyCollaborator struct {
	CollaboratorInfo `gorm:"embedded"`
	Account          CollaboratorAccount `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (CompanyCollaborator) TableName() string { return "collaborators" }
