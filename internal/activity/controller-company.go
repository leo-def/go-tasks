package activity

import (
	"net/http"

	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

type ControllerCompany struct {
	Service ServiceCompany
}

func NewControllerCompany(service ServiceCompany) *ControllerCompany {
	return &ControllerCompany{service}
}

// Get
// @Summary List activities
// @Description List activities in your company with pagination
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,title) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[CompanyActivityDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/ [get]
func (c *ControllerCompany) Get(ctx *gin.Context) {
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	params, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	items, count, err := c.Service.Get(params, companyID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	itemsDTO := ToCompanyActivityDTOs(items)
	payload := httpx.Paginated[CompanyActivityDTO]{Items: itemsDTO, Count: count}
	httpx.WriteOK(ctx, http.StatusOK, payload)
}

// GetById
// @Summary Get activity by ID
// @Description Get a single activity by ID
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Success 200 {object} httpx.Response[CompanyActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/{id} [get]
func (c *ControllerCompany) GetById(ctx *gin.Context) {
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	activity, found, err := c.Service.GetById(idUUID, companyID)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "activity not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToCompanyActivityDTO(activity))
}

// Update
// @Summary Update activity
// @Description Update an existing activity
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Param activity body ActivityUpdateDTO true "Activity update payload"
// @Success 200 {object} httpx.Response[CompanyActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/{id} [put]
func (c *ControllerCompany) Update(ctx *gin.Context) {
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto ActivityUpdateDTO
	if err = ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	activity := CompanyActivityFromUpdateDTOWithID(idUUID, dto)
	if err = c.Service.Update(&activity, companyID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, ToCompanyActivityDTO(&activity))
}

// Delete
// @Summary Delete activity
// @Description Delete activity by ID
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/{id} [delete]
func (c *ControllerCompany) Delete(ctx *gin.Context) {
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if found, err := c.Service.Delete(idUUID, companyID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "activity not found", nil)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "activity deleted: " + idUUID.String()})
}

// UpdateStatus
// @Summary Update activity status
// @Description Update status of an activity
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Activity ID" format(uuid)
// @Param status body ActivityUpdateStatusDTO true "Status update payload"
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/{id}/status [patch]
func (c *ControllerCompany) UpdateStatus(ctx *gin.Context) {
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto ActivityUpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := c.Service.UpdatedStatus(idUUID, lifecycle.LifecycleStatus(dto.Status), companyID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "activity status updated"})
}

// Create
// @Summary Create activity
// @Description Create a new activity
// @Tags Activity | Company Manager or Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param activity body ActivityCreateDTO true "Activity create payload"
// @Success 201 {object} httpx.Response[CompanyActivityDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/manager-owner/activity/ [post]
func (c *ControllerCompany) Create(ctx *gin.Context) {
	createdBy, found := jwttoken.GetCollaboratorID(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusBadRequest, "collaborator not found", nil)
		return
	}
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto ActivityCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	activity := CompanyActivityFromCreateDTO(dto)
	if err := c.Service.Create(&activity, createdBy, companyID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpx.WriteOK(ctx, http.StatusCreated, ToCompanyActivityDTO(&activity))
}
