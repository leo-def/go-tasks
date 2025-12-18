package auth

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdatePassword updates a user's password
// @Summary Update a user's password
// @Description Update the password of a user with the provided ID
// @Tags Auth | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param password body UpdatePasswordDTO true "Password data"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/auth/password [put]
func (c *Controller) UpdatePassword(ctx *gin.Context) {
	authData, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	accountID := authData.Id
	var dto UpdatePasswordDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if err := c.Service.UpdatePassword(accountID, dto.Password); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "password updated"})
}

// UpdateEmail updates a user's email
// @Summary Update a user's email
// @Description Update the email of a user with the provided ID
// @Tags Auth | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param email body UpdateEmailDTO true "Email data"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/auth/email [put]
func (c *Controller) UpdateEmail(ctx *gin.Context) {
	authData, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	accountID := authData.Id
	var dto UpdateEmailDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if err := c.Service.UpdateEmail(accountID, dto.Email); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "email updated"})
}

// UpdatePhone updates a user's phone
// @Summary Update a user's phone
// @Description Update the phone of a user with the provided ID
// @Tags Auth | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param phone body UpdatePhoneDTO true "Phone data"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/auth/phone [put]
func (c *Controller) UpdatePhone(ctx *gin.Context) {
	authData, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	accountID := authData.Id
	var dto UpdatePhoneDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if err := c.Service.UpdatePhone(accountID, dto.Phone); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "phone updated"})
}

// GetMe returns the authenticated user's information
// @Summary Get authenticated user's information
// @Description Get the information of the authenticated user
// @Tags Auth | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} httpx.Response[AuthDataDTO]
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/auth/me [get]
func (c *Controller) GetMe(ctx *gin.Context) {
	authData, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToAuthDataDTO(*authData))
}

// LoadCollaboratorContext loads the context of a collaborator
// @Summary Load collaborator context
// @Description Load the context of a collaborator with the provided ID
// @Tags Auth | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator_id body LoadCollaboratorContextDTO true "Collaborator ID"
// @Success 200 {object} httpx.Response[AuthTokenResposne]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/auth/collaborator [post]
func (c *Controller) LoadCollaboratorContext(ctx *gin.Context) {
	authData, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto LoadCollaboratorContextDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	token, err := c.Service.LoadCollaboratorContext(dto.CollaboratorID, *authData)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToAuthTokenResponse(token))

}
