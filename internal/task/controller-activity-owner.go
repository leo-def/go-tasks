package task

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerActivityOwner struct {
	Service ServiceActivityOwner
}

func NewControllerActivityOwner(service ServiceActivityOwner) *ControllerActivityOwner {
	return &ControllerActivityOwner{service}
}

// Create
// @Summary Create own task
// @Description Create a new task within an activity owned by the current collaborator
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param task body TaskCreateDTO true "Task create payload"
// @Success 201 {object} httpx.Response[TaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/ [post]
func (c *ControllerActivityOwner) Create(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto TaskCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	task := FromCreateDTO(dto)
	if err := c.Service.Create(&task, activityID, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusCreated, ToTaskDTO(&task))
}

// Delete
// @Summary Delete own task
// @Description Delete a task by ID within an activity owned by the current collaborator
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/{id} [delete]
func (c *ControllerActivityOwner) Delete(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	ok, err := c.Service.Delete(id, activityID, ownerID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if !ok {
		httpx.WriteError(ctx, http.StatusNotFound, "task not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "task deleted: " + id.String()})
}

// Get
// @Summary List own tasks by activity
// @Description List tasks by activity ID owned by the current collaborator with pagination
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,title) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[ActivityTaskDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/ [get]
func (c *ControllerActivityOwner) Get(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	params, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	items, count, err := c.Service.Get(params, activityID, ownerID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[ActivityTaskDTO]{Items: ToActivityTaskDTOs(items), Count: count})
}

// GetById
// @Summary Get own task
// @Description Get a single task by ID within an activity owned by the current collaborator
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} httpx.Response[ActivityTaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/{id} [get]
func (c *ControllerActivityOwner) GetById(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	item, ok, err := c.Service.GetById(id, activityID, ownerID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if !ok {
		httpx.WriteError(ctx, http.StatusNotFound, "task not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToActivityTaskDTO(item))
}

// Update
// @Summary Update own task
// @Description Update a task by ID within an activity owned by the current collaborator
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Param task body TaskUpdateDTO true "Task update payload"
// @Success 200 {object} httpx.Response[TaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/{id} [put]
func (c *ControllerActivityOwner) Update(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto TaskUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	model := FromUpdateDTOWithId(dto, id)
	if err := c.Service.Update(&model, activityID, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToTaskDTO(&model))
}

// UpdateStatus
// @Summary Update own task status
// @Description Update the status of a task within an activity owned by the current collaborator
// @Tags Task | Company Own Activity
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Param status body TaskUpdateStatusDTO true "Status update payload"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/activity/own/tasks/{activity_id}/{id}/status [patch]
func (c *ControllerActivityOwner) UpdateStatus(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	ownerID, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto TaskUpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := c.Service.UpdateStatus(id, lifecycle.LifecycleStatus(dto.Status), activityID, ownerID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "task status updated"})
}
