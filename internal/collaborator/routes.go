package collaborator

import "github.com/gin-gonic/gin"

func (m *Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	group := r.Group("/collaborator")
	{
		group.GET("/my-roles", m.Controller.MyRoles)
	}
}

func (m *Module) RegisterCompanyOwnerRoutes(r *gin.RouterGroup) {
	group := r.Group("/collaborator")
	{
		group.GET("/", m.ControllerCompany.Get)
		group.GET("/:id", m.ControllerCompany.GetById)
		group.POST("/", m.ControllerCompany.Create)
		group.POST("/with-account", m.ControllerCompany.CreateWithAccount)
		group.PUT("/:id", m.ControllerCompany.Update)
		group.DELETE("/:id", m.ControllerCompany.Delete)
	}
}

func (m *Module) RegisterCompanyRoutes(r *gin.RouterGroup) {
	group := r.Group("/collaborator/subordinates")
	{
		group.GET("/", m.ControllerSubCompany.Get)
		group.GET("/:id", m.ControllerSubCompany.GetById)
		group.POST("/", m.ControllerSubCompany.Create)
		group.POST("/with-account", m.ControllerSubCompany.CreateWithAccount)
		group.PUT("/:id", m.ControllerSubCompany.Update)
		group.DELETE("/:id", m.ControllerSubCompany.Delete)
	}
}
