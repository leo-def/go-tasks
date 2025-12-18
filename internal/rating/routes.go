package rating

import "github.com/gin-gonic/gin"

func (m *Module) RegisterCompanyManagerOwnerRoutes(r *gin.RouterGroup) {
    participationGroup := r.Group("/participation/rating/:participation_id")
    {
        participationGroup.POST("/", m.Controller.Create)
    }
    collaboratorGroup := r.Group("/collaborator/ratings/:collaborator_id")
    {
        collaboratorGroup.POST("/", m.Controller.CreateForCollaborator)
    }
}
