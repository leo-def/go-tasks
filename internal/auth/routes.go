package auth

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	group := r.Group("/auth")
	{
		group.POST("/sign-in", m.Controller.SignIn)
		group.POST("/sign-up", m.Controller.SignUp)
	}
}

func (m *Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	group := r.Group("/auth")
	{
		group.GET("/me", m.Controller.GetMe)
		group.PUT("/password", m.Controller.UpdatePassword)
		group.POST("/collaborator", m.Controller.LoadCollaboratorContext)
	}
}
