package jwttoken

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthData is the data structure that is used to store the authentication data in the JWT token.
type AuthData struct {
	Id           uuid.UUID
	Username     string
	Role         string
	SessionID    string
	Collaborator CollaboratorData
}

// CollaboratorData is the data structure that is used to store the collaborator data in the JWT token.
type CollaboratorData struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Role      string
}

// TokenData is the data structure that is used to store the token data in the JWT token.
type TokenData struct {
	jwt.RegisteredClaims
	Role             string
	CollaboratorID   string
	CompanyID        string
	CollaboratorRole string
}

func (a AuthData) HasID() bool             { return a.Id != uuid.Nil }
func (a AuthData) HasCollaboratorID() bool { return a.Collaborator.ID != uuid.Nil }
func (a AuthData) HasCompanyID() bool      { return a.Collaborator.CompanyID != uuid.Nil }

func (a AuthData) HasRole() bool { return a.Role != "" }

func (c CollaboratorData) HasID() bool        { return c.ID != uuid.Nil }
func (c CollaboratorData) HasCompanyID() bool { return c.CompanyID != uuid.Nil }
func (c CollaboratorData) HasRole() bool      { return c.Role != "" }
