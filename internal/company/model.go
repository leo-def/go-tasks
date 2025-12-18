package company

import (
	"go-tasks/internal/account"
	"go-tasks/internal/collaborator"
	"time"

	"github.com/google/uuid"
)

type CompanyInfo struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `json:"title" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

func (CompanyInfo) TableName() string { return "companies" }

type Company struct {
	CompanyInfo `gorm:"embedded"`
	Owner       CompanyOwner `json:"owner" gorm:"foreignKey:CompanyID;references:ID"`
}

type CompanyWithOwner struct {
    CompanyInfo `gorm:"embedded"`
    // Owner is required when creating a company. Its account password is hashed by the account service.
    Owner CompanyOwner `json:"owner"`
}

func (CompanyWithOwner) TableName() string { return "companies" }

type CompanyOwner struct {
	ID        uuid.UUID                     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID uuid.UUID                     `json:"company_id" gorm:"type:uuid;not null"`
	AccountID uuid.UUID                     `json:"account_id" gorm:"type:uuid;not null"`
	Role      collaborator.CollaboratorRole `json:"role" gorm:"not null"`
	Active    bool                          `json:"active" gorm:"not null;default:true"`
	Account   CompanyOwnerAccount           `json:"account" gorm:"foreignKey:AccountID;references:ID"`
}

func (CompanyOwner) TableName() string { return "collaborators" }

type CompanyOwnerAccount struct {
	ID       uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string       `json:"username" gorm:"unique;not null"`
	Password string       `json:"password" gorm:"not null"`
	Name     string       `json:"name" gorm:"not null"`
	Phone    string       `json:"phone" gorm:"not null"`
	Email    string       `json:"email" gorm:"unique;not null"`
	Role     account.Role `json:"role" gorm:"not null"`
}

func (CompanyOwnerAccount) TableName() string { return "accounts" }
