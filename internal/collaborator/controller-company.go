package collaborator

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ControllerCompany struct{ Service ServiceCompany }

func NewControllerCompany(service ServiceCompany) *ControllerCompany {
	return &ControllerCompany{service}
}

// Get
// @Summary List collaborators of my company
// @Description List collaborators in the authenticated user's company with pagination
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param skip query int false "Number of items to skip" default(0)
// @Param limit query int false "Maximum number of items to return" default(20)
// @Param sortBy query string false "Field to sort by" Enums(id,role) default(id)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC) default(ASC)
// @Param filter query string false "Filter expression"
// @Success 200 {object} httpx.Response[httpx.Paginated[CompanyCollaboratorDTO]]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/ [get]
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
	if collaborators, count, err := c.Service.Get(params, companyID); err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		dtos := ToCompanyCollaboratorDTOs(collaborators)
		httpx.WriteOK(ctx, http.StatusOK, httpx.Paginated[CompanyCollaboratorDTO]{Items: dtos, Count: count})
	}
}

// GetById
// @Summary Get collaborator by ID
// @Description Get a company collaborator by ID
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Success 200 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/{id} [get]
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
	if item, ok, err := c.Service.GetById(idUUID, companyID); err != nil {
		if err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else if !ok {
		httpx.WriteError(ctx, http.StatusNotFound, "not found", nil)
		return
	} else {
		dto := ToCompanyCollaboratorDTO(*item)
		httpx.WriteOK(ctx, http.StatusOK, dto)
	}
}

// Delete
// @Summary Delete collaborator
// @Description Delete a company collaborator by ID
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/{id} [delete]
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
	if ok, err := c.Service.Delete(idUUID, companyID); err != nil {
		if err == ErrForbiddenNotInCompany {
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

// Update
// @Summary Update collaborator
// @Description Update a company collaborator by ID
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Collaborator ID" format(uuid)
// @Param collaborator body CollaboratorUpdateDTO true "Collaborator update payload"
// @Success 200 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/{id} [put]
func (c *ControllerCompany) Update(ctx *gin.Context) {
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
	if err = c.Service.Update(&model, companyID); err != nil {
		if err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusOK, ToCompanyCollaboratorDTO(model))
	}
}

// Create
// @Summary Create collaborator
// @Description Create a new collaborator linked to an existing account
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator body CollaboratorCreateDTO true "Collaborator create payload"
// @Success 201 {object} httpx.Response[CollaboratorInfoDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/ [post]
func (c *ControllerCompany) Create(ctx *gin.Context) {
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
	info := FromCreateDTOWithCompanyID(dto, companyID)
	if err := c.Service.Create(&info, companyID); err != nil {
		if err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToCollaboratorInfoDTO(info))
	}
}

// CreateWithAccount
// @Summary Create collaborator with new account
// @Description Create a new collaborator and account within the company
// @Tags Collaborator | Company Owner
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param collaborator body CollaboratorWithAccountCreateDTO true "Collaborator with account create payload"
// @Success 201 {object} httpx.Response[CompanyCollaboratorDTO]
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /company/owner/collaborator/with-account [post]
func (c *ControllerCompany) CreateWithAccount(ctx *gin.Context) {
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
	if err := c.Service.CreateWithAccount(&model, dto.Account.Password, model.ID); err != nil {
		if err == ErrForbiddenNotInCompany {
			httpx.WriteError(ctx, http.StatusForbidden, err.Error(), err)
			return
		}
		httpx.WriteError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	} else {
		httpx.WriteOK(ctx, http.StatusCreated, ToCompanyCollaboratorDTO(model))
	}
}
