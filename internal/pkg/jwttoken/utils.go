package jwttoken

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthData(c *gin.Context) (*AuthData, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	authData, ok := user.(*AuthData)
	return authData, ok
}

func GetCompanyId(c *gin.Context) (uuid.UUID, bool) {
	authData, found := GetAuthData(c)
	if !found {
		return uuid.Nil, false
	}
	if !authData.HasCompanyID() {
		return uuid.Nil, false
	}
	return authData.Collaborator.CompanyID, true
}

func GetCollaboratorID(c *gin.Context) (uuid.UUID, bool) {
	authData, found := GetAuthData(c)
	if !found {
		return uuid.Nil, false
	}
	if !authData.HasCollaboratorID() {
		return uuid.Nil, false
	}
	return authData.Collaborator.ID, true
}

func GetRole(c *gin.Context) (string, bool) {
	authData, found := GetAuthData(c)
	if !found {
		return "", false
	}
	if !authData.HasRole() {
		return "", false
	}
	return authData.Collaborator.Role, true
}
