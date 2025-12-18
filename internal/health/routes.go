package health

import "github.com/gin-gonic/gin"

func (m *Module) RegisterRoutes(r *gin.Engine) {
    r.GET("/health", m.Controller.Get)
}

