package auth

import (
	"go-tasks/internal/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{ Service Service }

func NewController(service Service) *Controller {
	return &Controller{service}
}

// SignIn signs in a user and returns a JWT token
// @Summary Sign in a user
// @Description Sign in a user and return a JWT token
// @Tags Auth | Public
// @Accept json
// @Produce json
// @Param credentials body SignInDTO true "User credentials"
// @Success 200 {object} httpx.Response[AuthTokenResposne]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /auth/sign-in [post]
func (c *Controller) SignIn(ctx *gin.Context) {
	var dto SignInDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	token, err := c.Service.SignIn(dto.Username, dto.Password)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToAuthTokenResponse(token))

}

// SignUp signs up a new user
// @Summary Sign up a new user
// @Description Creates a new account with a securely hashed password via account service. Returns 204 with no body.
// @Tags Auth | Public
// @Accept json
// @Produce json
// @Param credentials body SignUpDTO true "User credentials"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /auth/sign-up [post]
func (c *Controller) SignUp(ctx *gin.Context) {
	var dto SignUpDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	account := FromSignUpDTO(dto)
	if err := c.Service.SignUp(&account, dto.Password); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	ctx.Status(http.StatusNoContent)
}
