package httpx

import (
	"go-tasks/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response is the standardized envelope for all endpoints
// data can be any type; error is an optional message for error responses
// path is the request path; timestamp is when the response is produced
//
// Use WriteOK and WriteError helpers to send this envelope.
type Response[T any] struct {
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
	Path      string    `json:"path"`
}

// WriteOK sends a standardized success response
func WriteOK[T any](ctx *gin.Context, status int, data T) {
	ctx.JSON(status, Response[T]{
		Data:      data,
		Timestamp: time.Now().UTC(),
		Error:     "",
		Path:      ctx.FullPath(),
	})
}

// WriteError sends a standardized error response
func WriteError(ctx *gin.Context, status int, errMsg string, err error) {
    if status >= http.StatusInternalServerError {
        fields := map[string]interface{}{"status": status}
        if err != nil {
            fields["error"] = err.Error()
        }
        logger.WithContext(ctx, "error", errMsg, fields)
        errMsg = "internal server error"
    }
    ctx.JSON(status, ErrorResponse{
        Data:      nil,
        Timestamp: time.Now().UTC(),
        Error:     errMsg,
        Path:      ctx.FullPath(),
    })
}

// Paginated is the common payload for list endpoints
// Items holds the current page; Count is the total matching records
// Use with WriteOK to wrap inside the Response envelope
type Paginated[T any] struct {
	Items []T   `json:"items"`
	Count int64 `json:"count"`
}

type ErrorResponse struct {
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
	Path      string    `json:"path"`
}

type MessageResponse struct {
	Data      MessageDTO `json:"data"`
	Timestamp time.Time  `json:"timestamp"`
	Error     string     `json:"error"`
	Path      string     `json:"path"`
}

type MessageDTO struct {
	Message string `json:"message" example:"operation completed"`
}
