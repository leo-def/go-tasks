package activity

import (
	"net/http"

	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

type ControllerOwn struct {
	Service ServiceOwn
}

func NewControllerOwn(service ServiceOwn) *ControllerOwn {
	return &ControllerOwn{service}
}

// Create
// @Summary Create activity
// @Description Create a new activity
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity body OwnActivityCreateDTO true "Activity create payload"
// @Success 201 {object} httpx.Response[OwnActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/ [post]
func (c *ControllerOwn) Create(ctx *gin.Context) {
	createdBy, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	var dto OwnActivityCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	activity := OwnActivityFromCreateDTO(dto)
	if err := c.Service.Create(&activity, createdBy, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusCreated, ToOwnActivityDTO(&activity))
}

// Delete
// @Summary Delete activity
// @Description Delete activity by ID
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/{id} [delete]
func (c *ControllerOwn) Delete(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	if found, err := c.Service.Delete(idUUID, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "activity not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "activity deleted: " + idUUID.String()})
}

// Get
// @Summary List activities
// @Description List your activities with pagination
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,title) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[OwnActivityDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/ [get]
func (c *ControllerOwn) Get(ctx *gin.Context) {
	params, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	items, count, err := c.Service.Get(params, ownerID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	itemsDTO := ToOwnActivityDTOs(items)
	payload := httpx.Paginated[OwnActivityDTO]{Items: itemsDTO, Count: count}
	httpx.WriteOK(ctx, http.StatusOK, payload)
}

// GetById
// @Summary Get owned activity by ID
// @Description Get a single activity owned by current user
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Success 200 {object} httpx.Response[OwnActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/{id} [get]
func (c *ControllerOwn) GetById(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	activity, found, err := c.Service.GetById(idUUID, ownerID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "activity not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToOwnActivityDTO(activity))
}

// Update
// @Summary Update activity
// @Description Update an existing activity
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Param activity body ActivityUpdateDTO true "Activity update payload"
// @Success 200 {object} httpx.Response[OwnActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/{id} [put]
func (c *ControllerOwn) Update(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	var dto OwnActivityUpdateDTO
	if err = ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	activity := OwnActivityFromUpdateDTOWithID(idUUID, dto)
	if err = c.Service.Update(&activity, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToOwnActivityDTO(&activity))
}

// UpdateStatus
// @Summary Update activity status
// @Description Update status of an activity
// @Tags Activity | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Param status body ActivityUpdateStatusDTO true "Status update payload"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/{id}/status [patch]
func (c *ControllerOwn) UpdateStatus(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	var dto ActivityUpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := c.Service.UpdateStatus(idUUID, lifecycle.LifecycleStatus(dto.Status), ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "activity status updated"})
}
