package rating

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

// Create
// @Summary Create or update rating
// @Description Create a new rating or update existing rating for a participation by the authenticated collaborator. If a rating already exists from the same collaborator for this participation, it will be replaced.
// @Tags Rating | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param participation_id path string true "Participation ID" format(uuid)
// @Param rating body RatingCreateDTO true "Rating create payload"
// @Success 201 {object} httpx.Response[RatingDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/participation/rating/{participation_id}/ [post]
func (c *Controller) Create(ctx *gin.Context) {
	assigneeID, ok := jwttoken.GetCollaboratorID(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	participationID, err := httpx.ResolveUUIDParam(ctx, "participation_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto RatingCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	rating := FromCreateDTO(dto)
	if err := c.Service.CreateForParticipation(&rating, assigneeID, participationID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToRatingDTO(&rating))
	}
}

// CreateForCollaborator
// @Summary Create or update rating for collaborator
// @Description Create a new rating or update existing rating for a collaborator by the authenticated company manager or owner. If a rating already exists from the same collaborator for this activity participation, it will be replaced.
// @Tags Rating | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator_id path string true "Collaborator ID" format(uuid)
// @Param rating body ActivityRatingCreateDTO true "Rating create payload"
// @Success 201 {object} httpx.Response[RatingDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/collaborator/ratings/{collaborator_id}/ [post]
func (c *Controller) CreateForCollaborator(ctx *gin.Context) {
	assigneeID, ok := jwttoken.GetCollaboratorID(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	collaboratorID, err := httpx.ResolveUUIDParam(ctx, "collaborator_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto ActivityRatingCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	rating := FromActivityRatingCreateDTO(dto)
	if err := c.Service.CreateForCollaborator(&rating, assigneeID, dto.ActivityID, collaboratorID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToRatingDTO(&rating))
	}
}
