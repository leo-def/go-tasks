package participation

import "github.com/gin-gonic/gin"

func (m *Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	group := r.Group("/participation")
	{
		group.GET("/my", m.Controller.My)
	}
}

func (m *Module) RegisterCompanyManagerOwnerRoutes(r *gin.RouterGroup) {
	group := r.Group("/activity/participations/:activity_id")
	{
		group.GET("/", m.ControllerActivity.Get)
		group.POST("/", m.ControllerActivity.Create)
	}

}
