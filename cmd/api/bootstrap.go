package main

import (
	"fmt"
	"go-tasks/internal/app"
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/env"
	"go-tasks/internal/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	DatabaseDriver string
	DatabaseURL    string
	AppPort        string
	TrustedProxies []string
}

func LoadConfig() *Config {
	proxies := ParseCSV(env.GetEnv("TRUSTED_PROXIES", "127.0.0.1,172.16.0.0/12"))
	return &Config{
		DatabaseDriver: env.GetEnv("DATABASE_DRIVER", "sqlite"),
		DatabaseURL:    env.GetEnv("DATABASE_URL", ""),
		AppPort:        env.GetEnv("APP_PORT", "8080"),
		TrustedProxies: proxies,
	}
}

func ParseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

type Application struct {
	Router *gin.Engine
	Config *Config
	DB     database.Connection
}

func InitializeApp() *Application {
	cfg := LoadConfig()
	db := database.ConnectDB(cfg.DatabaseDriver, cfg.DatabaseURL)

	appModule := app.Initialize(db)

	appModule.SeedAdminIfNeeded()

	router := gin.New()
	router.Use(gin.Recovery())
	if env.GetEnv("HTTP_LOG", "false") == "true" {
		router.Use(gin.Logger())
	}
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		panic(err)
	}
	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			ginSwagger.PersistAuthorization(true),
		))
	appModule.RegisterRoutes(router)

	return &Application{Router: router, Config: cfg, DB: db}
}

func (app *Application) Run() {
	port := app.Config.AppPort
	if port == "" {
		port = "8080"
	}
	logger.Info("server: listening", map[string]interface{}{"port": port})
	if err := app.Router.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Error("server: failed to start", map[string]interface{}{"error": err.Error(), "port": port})
		panic(err)
	}
}
