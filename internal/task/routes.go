package task

import "github.com/gin-gonic/gin"

func (m *Module) RegisterCompanyManagerOwnerRoutes(r *gin.RouterGroup) {
	group := r.Group("/activity/tasks/:activity_id/")
	{
		group.GET("/", m.ControllerActivity.Get)
		group.GET("/:id", m.ControllerActivity.GetById)
		group.POST("/", m.ControllerActivity.Create)
		group.PUT("/:id", m.ControllerActivity.Update)
		group.DELETE("/:id", m.ControllerActivity.Delete)
		group.PATCH("/:id/status", m.ControllerActivity.UpdateStatus)
	}
}

func (m *Module) RegisterCompanyRoutes(r *gin.RouterGroup) {
    group := r.Group("/activity/tasks/:activity_id")
    {
        group.PATCH("/:id/perform", m.ControllerActivity.Perform)
    }
    owner := r.Group("/activity/own/tasks/:activity_id/")
    {
        owner.GET("/", m.ControllerActivityOwner.Get)
        owner.GET("/:id", m.ControllerActivityOwner.GetById)
        owner.POST("/", m.ControllerActivityOwner.Create)
        owner.PUT("/:id", m.ControllerActivityOwner.Update)
        owner.DELETE("/:id", m.ControllerActivityOwner.Delete)
        owner.PATCH("/:id/status", m.ControllerActivityOwner.UpdateStatus)
    }
}
