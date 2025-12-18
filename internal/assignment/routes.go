package assignment

import "github.com/gin-gonic/gin"

func (m *Module) RegisterCompanyManagerOwnerRoutes(base *gin.RouterGroup) {
	// Task-scoped assignment routes
	taskGroup := base.Group("/task/assignments/:task_id")
	{
		taskGroup.GET("/", m.Controller.Get)
		taskGroup.GET("/:id", m.Controller.GetById)
		taskGroup.POST("", m.Controller.Create)
		taskGroup.POST("/participation/:participation_id", m.Controller.CreateForParticipation)
		taskGroup.POST("/collaborator/:collaborator_id", m.Controller.CreateForCollaborator)
		taskGroup.DELETE("/:id", m.Controller.Delete)
		taskGroup.PATCH("/:id/deactivate", m.Controller.Deactivate)
	}
}
