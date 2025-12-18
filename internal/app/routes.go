package app

import (
	"go-tasks/internal/account"
	"go-tasks/internal/auth"
	"go-tasks/internal/collaborator"

	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(router *gin.Engine) {

	m.HealthModule.RegisterRoutes(router)

	api := router.Group("/api/v1")
	api.Use(auth.JWTAuthMiddleware(m.AuthModule.TokenService))
	m.AuthModule.RegisterRoutes(api)

	protectedApi := api.Group("/protected")
	protectedApi.Use(auth.RequireAuthMiddleware(nil, nil, false))
	m.AuthModule.RegisterProtectedRoutes(protectedApi)
	m.CollaboratorModule.RegisterProtectedRoutes(protectedApi)
	m.ParticipationModule.RegisterProtectedRoutes(protectedApi)

	adminOrOpsApi := api.Group("/admin-or-ops")
	adminOrOpsApi.Use(auth.RequireAuthMiddleware([]string{
		string(account.RoleAdmin),
		string(account.RoleOps),
	}, nil, false))
	m.CompanyModule.RegisterAdminOrOpsRoutes(adminOrOpsApi)

	companyApi := api.Group("/company")
	companyApi.Use(auth.RequireAuthMiddleware(nil, nil, true))
	m.ActivityModule.RegisterCompanyRoutes(companyApi)
	m.CollaboratorModule.RegisterCompanyRoutes(companyApi)
	m.TaskModule.RegisterCompanyRoutes(companyApi)

	companyOwnerApi := api.Group("/company/owner")
	companyOwnerApi.Use(auth.RequireAuthMiddleware(nil, []string{
		string(collaborator.CollaboratorRoleOwner),
	}, true))
	m.CollaboratorModule.RegisterCompanyOwnerRoutes(companyOwnerApi)

	companyManagerOwnerApi := api.Group("/company/manager-owner")
	companyManagerOwnerApi.Use(auth.RequireAuthMiddleware(nil, []string{
		string(collaborator.CollaboratorRoleOwner),
		string(collaborator.CollaboratorRoleManager),
	}, true))
	m.ActivityModule.RegisterCompanyManagerOwnerRoutes(companyManagerOwnerApi)
	m.AssignmentModule.RegisterCompanyManagerOwnerRoutes(companyManagerOwnerApi)
	m.ParticipationModule.RegisterCompanyManagerOwnerRoutes(companyManagerOwnerApi)
	m.RatingModule.RegisterCompanyManagerOwnerRoutes(companyManagerOwnerApi)
	m.TaskModule.RegisterCompanyManagerOwnerRoutes(companyManagerOwnerApi)

}
