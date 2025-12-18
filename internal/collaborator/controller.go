package collaborator

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{ Service Service }

func NewController(service Service) *Controller {
	return &Controller{service}
}

// MyRoles
// @Summary List my collaborator roles
// @Description List collaborator roles for the authenticated account
// @Tags Collaborator | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} httpx.Response[httpx.Paginated[AccountCollaboratorDTO]]
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/collaborator/my-roles [get]
func (c *Controller) MyRoles(ctx *gin.Context) {
	account, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	accountID := account.Id
	collaborators, count, err := c.Service.GetByAccountId(accountID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	dtos := ToAccountCollaboratorDTOs(collaborators)
	httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[AccountCollaboratorDTO]{Items: dtos, Count: count})
}
