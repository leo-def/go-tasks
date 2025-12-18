package activity

import "github.com/gin-gonic/gin"

func (m *Module) RegisterCompanyManagerOwnerRoutes(r *gin.RouterGroup) {
	group := r.Group("/activity")
	{
		group.GET("/", m.Controller.Get)
		group.GET("/:id", m.Controller.GetById)
		group.POST("/", m.Controller.Create)
		group.PUT("/:id", m.Controller.Update)
		group.DELETE("/:id", m.Controller.Delete)
		group.PATCH("/:id/status", m.Controller.UpdateStatus)
	}
}

func (m *Module) RegisterCompanyRoutes(r *gin.RouterGroup) {
	group := r.Group("/activity/own")
	{
		group.GET("/", m.ControllerOwn.Get)
		group.GET("/:id", m.ControllerOwn.GetById)
		group.POST("/", m.ControllerOwn.Create)
		group.PUT("/:id", m.ControllerOwn.Update)
		group.DELETE("/:id", m.ControllerOwn.Delete)
		group.PATCH("/:id/status", m.ControllerOwn.UpdateStatus)
	}
}
