package participation

import (
	"go-tasks/internal/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerActivity struct{ Service ServiceActivity }

func NewControllerActivity(service ServiceActivity) *ControllerActivity {
	return &ControllerActivity{service}
}

// Get
// @Summary List participations by activity
// @Description List participations for a given activity ID
// @Tags Participation | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.Paginated[ActivityParticipationDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/participations/{activity_id}/ [get]
func (c *ControllerActivity) Get(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if items, count, err := c.Service.Get(activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		dtos := ToActivityParticipationDTOs(items)
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[ActivityParticipationDTO]{Items: dtos, Count: count})
	}
}

// Create
// @Summary Create participation
// @Description Create a participation for an activity
// @Tags Participation | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param participation body ParticipationCreateDTO true "Participation create payload"
// @Success 201 {object} httpx.Response[ParticipationInfoDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/participations/{activity_id}/ [post]
func (c *ControllerActivity) Create(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto ParticipationCreateDTO
	if err = ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	participation := FromCreateDTO(dto)
	if err = c.Service.Create(&participation, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToParticipationInfoDTO(participation))
	}
}
