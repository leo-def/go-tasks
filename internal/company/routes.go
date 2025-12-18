package company

import "github.com/gin-gonic/gin"

func (m *Module) RegisterAdminOrOpsRoutes(r *gin.RouterGroup) {
	group := r.Group("/company")
	{
		group.GET("/", m.Controller.GetAll)
		group.GET("/:id", m.Controller.GetById)
		group.POST("/", m.Controller.Create)
		group.PUT("/:id", m.Controller.Update)
		group.DELETE("/:id", m.Controller.Delete)
	}
}
