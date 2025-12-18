package collaborator

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerSubCompany struct{ Service ServiceSubCompany }

func NewControllerSubCompany(service ServiceSubCompany) *ControllerSubCompany {
	return &ControllerSubCompany{Service: service}
}

// Get
// @Summary List sub-collaborators
// @Description List sub-collaborators filtered by allowed roles in my company
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,role) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[CompanyCollaborator]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/ [get]
func (c *ControllerSubCompany) Get(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
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
	if items, count, err := c.Service.GetByCompanyId(companyID, CollaboratorRole(role), params); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[CompanyCollaborator]{Items: items, Count: count})
	}
}

// GetById
// @Summary Get sub-collaborator by ID
// @Description Get a sub-collaborator and validate role access
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Success 200 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/{id} [get]
func (c *ControllerSubCompany) GetById(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
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
	if target, _, err := c.Service.GetById(idUUID, companyID, CollaboratorRole(role)); err != nil {
		if err == ErrForbiddenNotSubordinate || err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToCompanyCollaboratorDTO(*target))
	}
}

// Create
// @Summary Create sub-collaborator
// @Description Create a new collaborator with role validation
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator body CollaboratorCreateDTO true "Collaborator create payload"
// @Success 201 {object} httpx.Response[CollaboratorInfoDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/ [post]
func (c *ControllerSubCompany) Create(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto CollaboratorCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	model := FromCreateDTOWithCompanyID(dto, companyID)
	if err := c.Service.Create(&model, companyID, CollaboratorRole(role)); err != nil {
		if err == ErrForbiddenNotSubordinate || err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToCollaboratorInfoDTO(model))
	}
}

// Update
// @Summary Update sub-collaborator
// @Description Update an existing collaborator with role validation
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Param collaborator body CollaboratorUpdateDTO true "Collaborator update payload"
// @Success 200 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/{id} [put]
func (c *ControllerSubCompany) Update(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto CollaboratorUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	idUUID, err := httpx.ResolveUUIDParam(ctx, "id")
	if err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	model := CompanyCollaboratorFromUpdateDTOWithID(dto, idUUID, companyID)
	if err = c.Service.Update(&model, companyID, CollaboratorRole(role)); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToCompanyCollaboratorDTO(model))
	}
}

// Delete
// @Summary Delete sub-collaborator
// @Description Delete a collaborator with role validation
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/{id} [delete]
func (c *ControllerSubCompany) Delete(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
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
	if ok, err := c.Service.Delete(idUUID, companyID, CollaboratorRole(role)); err != nil {
		if err == ErrForbiddenNotSubordinate || err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !ok {
		httpx.WriteError(ctx, http.StatusNotFound, "not found", nil)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "collaborator deleted: " + idUUID.String()})
	}
}

// CreateWithAccount
// @Summary Create sub-collaborator with new account
// @Description Create a new sub-collaborator and account within the company
// @Tags Collaborator | Company
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator body CollaboratorWithAccountCreateDTO true "Collaborator with account create payload"
// @Success 201 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/collaborator/subordinates/with-account [post]
func (c *ControllerSubCompany) CreateWithAccount(ctx *gin.Context) {
	role, ok := jwttoken.GetRole(ctx)
	if !ok {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	companyID, found := jwttoken.GetCompanyId(ctx)
	if !found {
		httpx.WriteError(ctx, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	var dto CollaboratorWithAccountCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}
	model := CompanyCollaboratorFromCreateNewAccountDTO(dto, companyID)
	if err := c.Service.CreateWithAccount(&model, dto.Account.Password, companyID, CollaboratorRole(role)); err != nil {
		if err == ErrForbiddenNotSubordinate || err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToCompanyCollaboratorDTO(model))
	}
}
