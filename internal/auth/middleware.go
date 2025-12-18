package auth

import (
	"go-tasks/internal/pkg/httpx"
	"go-tasks/internal/pkg/jwttoken"
	"go-tasks/internal/pkg/logger"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(tokenService jwttoken.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
			token = strings.TrimSpace(token)

			authData, err := tokenService.ParseToken(token)
			if err != nil {
				logger.WithContext(c, "warn", "token parse failure", map[string]interface{}{"reason": err.Error()})
			} else {
				c.Set("user", authData)
			}
		} else {
			// Only warn, do not log token contents
			logger.WithContext(c, "warn", "no bearer token provided", nil)
		}
		c.Next()
	}
}

func RequireAuthMiddleware(roles []string, collaboratorRoles []string, companyAction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")

		if !exists {
			logger.WithContext(c, "error", "unauthorized: missing auth context", nil)
			httpx.WriteError(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}

		authData, ok := user.(*jwttoken.AuthData)
		if !ok {
			logger.WithContext(c, "error", "unauthorized: invalid auth context", nil)
			httpx.WriteError(c, http.StatusUnauthorized, "unauthorized", nil)
			c.Abort()
			return
		}
		if companyAction {
			if !authData.HasCompanyID() {
				logger.WithContext(c, "error", "unauthorized: missing company id", nil)
				httpx.WriteError(c, http.StatusUnauthorized, "unauthorized", nil)
				c.Abort()
				return
			}
		}
		if roles != nil {
			if !slices.Contains(roles, authData.Role) {
				logger.WithContext(c, "warn", "forbidden: role mismatch", map[string]interface{}{"required": roles, "actual": authData.Role})
				httpx.WriteError(c, http.StatusForbidden, "forbidden", nil)
				c.Abort()
				return
			}
		}
		if collaboratorRoles != nil {
			if !slices.Contains(collaboratorRoles, authData.Collaborator.Role) {
				logger.WithContext(c, "warn", "forbidden: company role mismatch", map[string]interface{}{"required": collaboratorRoles, "actual": authData.Collaborator.Role})
				httpx.WriteError(c, http.StatusForbidden, "forbidden", nil)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
