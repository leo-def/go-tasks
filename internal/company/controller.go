package company

import (
	"go-tasks/internal/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{ Service Service }

func NewController(service Service) *Controller {
	return &Controller{service}
}

// Create handles creating a new company
// @Summary Create a new company
// @Description Atomically creates a company, its owner account (hashed password via account service), and links the owner as a collaborator.
// @Tags Company | Admin or Ops
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param company body CompanyCreateDTO true "Company data"
// @Success 201 {object} httpx.Response[CompanyDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin-or-ops/company/ [post]
func (c *Controller) Create(ctx *gin.Context) {
	var dto CompanyCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	companyRequest := FromCreateDTO(dto)
	company, err := c.Service.Create(&companyRequest)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToCompanyDTO(company))
	}
}

// GetAll handles fetch with common pagination and response envelope
// @Summary Get all companies
// @Description Get a list of all companies with pagination
// @Tags Company | Admin or Ops
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,title) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[CompanyDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin-or-ops/company/ [get]
func (c *Controller) GetAll(ctx *gin.Context) {
	params, err := httpx.ResolvePagination(ctx)
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if items, count, err := c.Service.Get(params); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		itemsDTO := ToCompanyDTOs(items)
		payload := httpx.Paginated[CompanyDTO]{Items: itemsDTO, Count: count}
		httpx.WriteOK(ctx, http.StatusOK, payload)
	}
}

// GetById handles fetching a single company by ID
// @Summary Get company by ID
// @Description Get a company by its unique ID
// @Tags Company | Admin or Ops
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID" format(uuid)
// @Success 200 {object} httpx.Response[CompanyDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin-or-ops/company/{id} [get]
func (c *Controller) GetById(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if company, found, err := c.Service.GetById(idUUID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "company not found", nil)
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToCompanyDTO(*company))
	}
}

// Update handles updating an existing company
// @Summary Update a company
// @Description Update the details of an existing company
// @Tags Company | Admin or Ops
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID" format(uuid)
// @Param company body CompanyUpdateDTO true "Company data"
// @Success 200 {object} httpx.Response[CompanyDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin-or-ops/company/{id} [put]
func (c *Controller) Update(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	var dto CompanyUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	company := FromUpdateDTOWithID(idUUID, dto)
	if err := c.Service.Update(&company); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToCompanyDTO(company))
	}
}

// Delete handles deleting a company by ID
// @Summary Delete a company
// @Description Delete a company by its unique ID
// @Tags Company | Admin or Ops
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Company ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /admin-or-ops/company/{id} [delete]
func (c *Controller) Delete(ctx *gin.Context) {
	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	if found, err := c.Service.Delete(idUUID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
	} else if !found {
		httpx.WriteError(ctx, http.StatusNotFound, "company not found", nil)
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "company deleted: " + idUUID.String()})
	}

}
