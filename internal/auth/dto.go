package auth

import (
	"go-tasks/internal/account"
	"go-tasks/internal/pkg/jwttoken"

	"github.com/google/uuid"
)

type AuthRole string

const (
	AuthRoleUser  AuthRole = "USER"
	AuthRoleAdmin AuthRole = "ADMIN"
	AuthRoleOps   AuthRole = "OPS"
)

type AuthCollaboratorRole string

const (
	AuthCollaboratorRoleOwner   AuthCollaboratorRole = "OWNER"
	AuthCollaboratorRoleManager AuthCollaboratorRole = "MANAGER"
	AuthCollaboratorRoleOps     AuthCollaboratorRole = "OPS"
)

type SignInDTO struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

type SignUpDTO struct {
	// Password is hashed by account service.
	Username string `json:"username" binding:"required" example:"new_user"`
	Name     string `json:"name" binding:"required" example:"New User"`
	Phone    string `json:"phone" binding:"required" example:"+1-202-555-0170"`
	Email    string `json:"email" binding:"required" example:"new_user@example.com"`
	Password string `json:"password" binding:"required" example:"S3cureP@ss!"`
}

type UpdatePasswordDTO struct {
	Password string `json:"password" binding:"required"`
}

type UpdateEmailDTO struct {
	Email string `json:"email" binding:"required" example:"new_user@example.com"`
}

type UpdatePhoneDTO struct {
	Phone string `json:"phone" binding:"required" example:"+1-202-555-0170"`
}

type LoadCollaboratorContextDTO struct {
	CollaboratorID uuid.UUID `json:"collaborator_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type AuthTokenResposne struct {
	Token         string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	Authorization string `json:"authorization" example:"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

func FromSignUpDTO(dto SignUpDTO) account.Account {
	return account.Account{
		Username: dto.Username,
		Name:     dto.Name,
		Phone:    dto.Phone,
		Email:    dto.Email,
		Role:     account.RoleUser,
	}
}

func ToAuthTokenResponse(token string) AuthTokenResposne {
	return AuthTokenResposne{
		Token:         token,
		Authorization: "Bearer " + token,
	}
}

type AuthDataDTO struct {
	Id           string              `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username     string              `json:"username" example:"admin"`
	Role         AuthRole            `json:"role" example:"USER"`
	SessionID    string              `json:"sessionID" example:"sess-abc-123"`
	Collaborator CollaboratorDataDTO `json:"collaborator"`
}

type CollaboratorDataDTO struct {
	ID        string               `json:"id" example:"223e4567-e89b-12d3-a456-426614174000"`
	CompanyID string               `json:"companyID" example:"323e4567-e89b-12d3-a456-426614174000"`
	Role      AuthCollaboratorRole `json:"role" example:"OWNER"`
}

func ToAuthDataDTO(a jwttoken.AuthData) AuthDataDTO {
	return AuthDataDTO{
		Id:        a.Id.String(),
		Username:  a.Username,
		Role:      AuthRole(a.Role),
		SessionID: a.SessionID,
		Collaborator: CollaboratorDataDTO{
			ID:        a.Collaborator.ID.String(),
			CompanyID: a.Collaborator.CompanyID.String(),
			Role:      AuthCollaboratorRole(a.Collaborator.Role),
		},
	}
}
