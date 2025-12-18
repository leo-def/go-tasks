package task

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerActivity struct {
	Service ServiceActivity
}

func NewControllerActivity(service ServiceActivity) *ControllerActivity {
	return &ControllerActivity{service}
}

// Create
// @Summary Create task
// @Description Create a new task within an activity
// @Tags Task | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param task body TaskCreateDTO true "Task create payload"
// @Success 201 {object} httpx.Response[TaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/ [post]
func (c *ControllerActivity) Create(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	var dto TaskCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	task := FromCreateDTO(dto)
	if err := c.Service.Create(&task, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpx.WriteOK(ctx, http.StatusCreated, ToTaskDTO(&task))
}

// Delete
// @Summary Delete task
// @Description Delete a task by ID within an activity
// @Tags Task | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/{id} [delete]
func (c *ControllerActivity) Delete(ctx *gin.Context) {
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

	if found, err := c.Service.Delete(id, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "task not found", nil)
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "task deleted: " + id.String()})
	}

}

// Get
// @Summary List tasks by activity
// @Description List tasks by activity ID with pagination
// @Tags Task | Company Manager or Owner
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
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/ [get]
func (c *ControllerActivity) Get(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	params, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	if tasks, total, err := c.Service.Get(params, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		tasksDTO := ToActivityTaskDTOs(tasks)
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[ActivityTaskDTO]{
			Items: tasksDTO,
			Count: total,
		})
	}
}

// GetById
// @Summary Get task
// @Description Get a single task by ID within an activity
// @Tags Task | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} httpx.Response[ActivityTaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/{id} [get]
func (c *ControllerActivity) GetById(ctx *gin.Context) {
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

	if item, found, err := c.Service.GetById(id, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "task not found", nil)
		return
	} else {
		dto := ToActivityTaskDTO(item)
		httpx.WriteOK(ctx, http.StatusOK, dto)
	}
}

// Perform
// @Summary Perform task
// @Description Perform a task by updating its status for current collaborator
// @Tags Task | Company Manager or Owner
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
// @Router /company/activity/tasks/{activity_id}/{id}/perform [patch]
func (c *ControllerActivity) Perform(ctx *gin.Context) {
	activityID, err := httpx.ResolveUUIDParam(ctx, "activity_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	auth, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	collaboratorID := auth.Collaborator.ID

	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	var dto TaskUpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := c.Service.Perform(id, collaboratorID, lifecycle.LifecycleStatus(dto.Status), activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "task status updated"})

}

// Update
// @Summary Update task
// @Description Update a task by ID within an activity
// @Tags Task | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Param task body TaskUpdateDTO true "Task update payload"
// @Success 200 {object} httpx.Response[TaskDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/{id} [put]
func (c *ControllerActivity) Update(ctx *gin.Context) {
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

	var dto TaskUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	model := FromUpdateDTOWithId(dto, id)
	if err := c.Service.Update(&model, activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpx.WriteOK(ctx, http.StatusOK, ToTaskDTO(&model))
}

// UpdateStatus
// @Summary Update task status
// @Description Update the status of a task within an activity
// @Tags Task | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity_id path string true "Activity ID" format(uuid)
// @Param id path string true "Task ID" format(uuid)
// @Param status body TaskUpdateStatusDTO true "Status update payload"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/tasks/{activity_id}/{id}/status [patch]
func (c *ControllerActivity) UpdateStatus(ctx *gin.Context) {
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

	var dto TaskUpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := c.Service.UpdateStatus(id, lifecycle.LifecycleStatus(dto.Status), activityID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "task status updated"})
}
