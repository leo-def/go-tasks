package assignment

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ControllerTask struct{ Service ServiceTask }

func NewControllerTask(service ServiceTask) *ControllerTask {
	return &ControllerTask{Service: service}
}

// Get
// @Summary List assignments by task
// @Description List all assignments (active and inactive) for a given task ID with pagination support. Returns assignment history including deactivated assignments to provide full audit trail.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,assign_date) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[TaskAssignmentDTO]] "List of assignments with pagination metadata"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error"
// @Router /company/manager-owner/task/assignments/{task_id}/ [get]
func (c *ControllerTask) Get(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	pagination, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if assignments, count, err := c.Service.GetByTaskId(pagination, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		dtos := ToTaskAssignmentDTOs(assignments)
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[TaskAssignmentDTO]{Items: dtos, Count: count})
	}
}

// GetById
// @Summary Get assignment by ID
// @Description Get a specific assignment by ID within a task, including full details of assigned collaborator, assigner, and task information. Returns complete assignment details with nested objects.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param id path string true "Assignment ID" format(uuid)
// @Success 200 {object} httpx.Response[AssignmentDTO] "Assignment details with full nested information"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 404 {object} httpx.ErrorResponse "Assignment not found"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error"
// @Router /company/manager-owner/task/assignments/{task_id}/{id} [get]
func (c *ControllerTask) GetById(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, "invalid id", err)
		return
	}
	if assignment, found, err := c.Service.GetById(id, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, "failed to get assignment", err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "assignment not found", nil)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToAssignmentDTO(*assignment))
	}
}

// Delete
// @Summary Delete assignment
// @Description Permanently delete an assignment by ID. This physically removes the assignment from the database. **Warning**: This action cannot be undone. Consider using deactivate instead to maintain audit trail.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param id path string true "Assignment ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO] "Assignment deleted successfully"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 404 {object} httpx.ErrorResponse "Assignment not found"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error"
// @Router /company/manager-owner/task/assignments/{task_id}/{id} [delete]
func (c *ControllerTask) Delete(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if success, err := c.Service.Delete(id, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !success {
		httpx.WriteError(ctx, http.StatusNotFound, "assignment not found", nil)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "assignment deleted: " + id.String()})
	}
}

// Deactivate
// @Summary Deactivate assignment
// @Description Manually deactivate an assignment by setting its active status to false. The assignment remains in the database but is no longer considered active. This is the recommended way to end an assignment while maintaining audit trail.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param id path string true "Assignment ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO] "Assignment deactivated successfully"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 404 {object} httpx.ErrorResponse "Assignment not found"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error"
// @Router /company/manager-owner/task/assignments/{task_id}/{id}/deactivate [patch]
func (c *ControllerTask) Deactivate(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	id, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	if found, err := c.Service.Deactivate(id, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "assignment not found", nil)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "assignment deactivated: " + id.String()})
	}
}

// Create
// @Summary Create assignment (direct method)
// @Description Create a new assignment for a task using a direct participation ID. This method requires you to provide the participation ID directly in the request body. **Business Rule**: If there are existing active assignments for this task, they will be automatically deactivated to ensure only one active assignment per task. The operation is performed within a database transaction to ensure consistency.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param assignment body AssignmentCreateDTO true "Assignment create payload with participation ID"
// @Success 201 {object} httpx.Response[AssignmentDTO] "Assignment created successfully"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters or JSON payload"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error or database transaction failed"
// @Router /company/manager-owner/task/assignments/{task_id} [post]
func (c *ControllerTask) Create(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	auth, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	assignerID := auth.Collaborator.ID
	var dto AssignmentCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	assignment := FromAssignmentCreateDTO(dto)
	if err := c.Service.Create(&assignment, assignerID, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToAssignmentDTO(assignment))
	}
}

// CreateForParticipation
// @Summary Create assignment for participation
// @Description Create a new assignment for a task using participation ID from the URL path. This method provides a cleaner API by taking the participation ID from the URL rather than the request body. **Business Rule**: If there are existing active assignments for this task, they will be automatically deactivated to ensure only one active assignment per task. The operation is performed within a database transaction to ensure consistency.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param participation_id path string true "Participation ID (collaborator + activity relationship)" format(uuid)
// @Param assignment body AssignmentParticipationCreateDTO true "Assignment create payload (no participation ID needed)"
// @Success 201 {object} httpx.Response[AssignmentDTO] "Assignment created successfully"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters or JSON payload"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 404 {object} httpx.ErrorResponse "Participation not found"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error or database transaction failed"
// @Router /company/manager-owner/task/assignments/{task_id}/participation/{participation_id} [post]
func (c *ControllerTask) CreateForParticipation(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	participationID, err := httpx.ResolveUUIDParam(ctx, "participation_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	auth, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	assignerID := auth.Collaborator.ID
	var dto AssignmentParticipationCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	assignment := FromAssignmentParticipationCreateDTO(dto)
	if err := c.Service.CreateForParticipation(&assignment, assignerID, participationID, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToAssignmentDTO(assignment))
	}
}

// CreateForCollaborator
// @Summary Create assignment for collaborator
// @Description Create a new assignment for a task using collaborator ID and activity ID. The system will automatically find an existing participation or create a new one for the collaborator+activity combination. This is the most convenient method when you have collaborator and activity IDs. **Business Rule**: If there are existing active assignments for this task, they will be automatically deactivated to ensure only one active assignment per task. The operation is performed within a database transaction to ensure consistency.
// @Tags Assignment | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Param collaborator_id path string true "Collaborator ID (person being assigned)" format(uuid)
// @Param assignment body AssignmentCollaboratorCreateDTO true "Assignment create payload with activity ID"
// @Success 201 {object} httpx.Response[AssignmentDTO] "Assignment created successfully (participation auto-created if needed)"
// @Failure 400 {object} httpx.ErrorResponse "Invalid request parameters or JSON payload"
// @Failure 401 {object} httpx.ErrorResponse "Authentication required"
// @Failure 404 {object} httpx.ErrorResponse "Collaborator or activity not found"
// @Failure 500 {object} httpx.ErrorResponse "Internal server error or database transaction failed"
// @Router /company/manager-owner/task/assignments/{task_id}/collaborator/{collaborator_id} [post]
func (c *ControllerTask) CreateForCollaborator(ctx *gin.Context) {
	taskID, err := httpx.ResolveUUIDParam(ctx, "task_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	collaboratorID, err := httpx.ResolveUUIDParam(ctx, "collaborator_id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	auth, found := jwttoken.GetAuthData(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	assignerID := auth.Collaborator.ID
	var dto AssignmentCollaboratorCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	assignment := FromAssignmentCollaboratorCreateDTO(dto)
	if err := c.Service.CreateForCollaborator(&assignment, assignerID, dto.ActivityID, collaboratorID, taskID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToAssignmentDTO(assignment))
	}
}
