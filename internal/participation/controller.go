package participation

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

// My
// @Summary List my participations
// @Description List participations for the authenticated collaborator
// @Tags Participation | Protected
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} httpx.Response[httpx.Paginated[CollaboratorParticipationDTO]]
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /protected/participation/my [get]
func (c *Controller) My(ctx *gin.Context) {
	auth, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	collaboratorID := auth.Collaborator.ID
	if items, count, err := c.Service.GetByCollaboratorId(collaboratorID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		dtos := ToCollaboratorParticipationDTOs(items)
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[CollaboratorParticipationDTO]{Items: dtos, Count: count})
	}
}
