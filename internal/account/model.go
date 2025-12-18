package account

import (
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents the role of an account as an enum-like type
type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleOps   Role = "OPS"
	RoleUser  Role = "USER"
)

type Account struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      Role      `json:"role" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

type AuthAccount struct {
	Account `gorm:"embedded"`
	Roles   []AccountCollaborator `gorm:"foreignKey:AccountID;references:ID"`
}

func (AuthAccount) TableName() string {
	return "accounts"
}

type AccountCollaborator struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CompanyID uuid.UUID `json:"company_id" gorm:"type:uuid;not null"`
	AccountID uuid.UUID `json:"account_id" gorm:"type:uuid;not null"`
	Role      string    `json:"role" gorm:"not null"`
	Active    bool      `json:"active" gorm:"not null;default:true"`
}

func (AccountCollaborator) TableName() string {
	return "collaborators"
}

// BeforeSave is a GORM hook that sanitizes the account's username before saving it to the database.
// It trims leading and trailing whitespace and escapes HTML characters to prevent XSS attacks.
func (account *Account) BeforeSave(tx *gorm.DB) error {
	account.Username = html.EscapeString(strings.TrimSpace((account.Username)))
	return nil
}
