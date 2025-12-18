package app

import (
	"go-tasks/internal/pkg/env"
	"go-tasks/internal/pkg/logger"
	"strings"
)

func (m *Module) SeedAdminIfNeeded() {
	appEnv := strings.ToLower(env.GetEnv("APP_ENV", "local"))
	if appEnv == "prod" || appEnv == "production" {
		return
	}
	username := env.GetEnv("ADMIN_USERNAME", "admin")
	password := env.GetEnv("ADMIN_PASSWORD", "admin123")
	email := env.GetEnv("ADMIN_EMAIL", "admin@example.com")
	name := env.GetEnv("ADMIN_NAME", "Admin User")
	phone := env.GetEnv("ADMIN_PHONE", "+5511999999999")

	if err := (*m.AccountModule.Service).EnsureAdmin(username, password, email, name, phone); err != nil {
		logger.Error("seed: failed to ensure admin", map[string]interface{}{"reason": err.Error()})
		return
	}
	logger.Info("seed: admin user ensured", map[string]interface{}{"username": username})
}
