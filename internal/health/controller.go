package health

import (
	"go-tasks/internal/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func NewController() *Controller { return &Controller{} }

// Get
// @Summary Health check
// @Description Check if the service is up and running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} httpx.Response[httpx.MessageDTO]
// @Router /health [get]
func (c *Controller) Get(ctx *gin.Context) {
	httpx.WriteOK(ctx, http.StatusOK, httpx.MessageDTO{Message: "OK"})
}
