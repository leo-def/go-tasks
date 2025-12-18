# Project Structure - Modular Go Application (Gin + Gorm)

This document describes the recommended structure for a modular Go application using:
- **Gin** (HTTP router)
- **GORM** (ORM)
- **Dependency Injection-like initialization**
- **Environment-based configuration**
- **Isolated modules with their own routes, models, services, repositories**

---

## 📁 Folder Structure

```
/your-project
├── cmd/
│   └── api/
│       └── main.go          # Entry point: loads config, modules, server
├── config/
│   ├── config.go            # Reads environment variables & loads .env
│   └── database.go          # DB initialization using GORM
├── internal/
│   ├── server/
│   │   └── server.go        # Creates Gin engine and loads modules
│   ├── core/
│   │   ├── module.go        # Interface for modules (Init, RegisterRoutes)
│   │   └── container.go     # Root container for dependencies
│   ├── user/                # Example module: User
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── module.go        # Registers repository, service, handler
│   └── product/             # Example module: Product (same structure)
├── pkg/                     # Utilities shared across modules
│   ├── logger/
│   └── response/
├── .env                     # Environment variables (local only)
├── go.mod
├── go.sum
└── README.md
```

---

## ✅ Module Pattern

Each module (like `user`, `product`) has its own:

| File        | Responsibility                          |
|-------------|-------------------------------------------|
| `model.go`  | Struct definitions (GORM models)         |
| `repository.go` | DB operations (CRUD)                |
| `service.go` | Business logic                         |
| `handler.go` | HTTP routing handlers (Gin)            |
| `module.go`  | Setup: creates repository, service, handler |

---

## 🛠 Example: Module Registration (`user/module.go`)

```go
package user

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Module struct {
    Repo    Repository
    Service Service
    Handler Handler
}

func NewModule(db *gorm.DB) *Module {
    repo := NewRepository(db)
    service := NewService(repo)
    handler := NewHandler(service)
    return &Module{repo, service, handler}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
    group := router.Group("/users")
    group.GET("/", m.Handler.GetAll)
    group.POST("/", m.Handler.Create)
}
```

---

## 🧠 Server Bootstrapping (`cmd/api/main.go`)

```go
package main

import (
    "your-project/config"
    "your-project/internal/server"
)

func main() {
    config.LoadEnv()
    db := config.ConnectDB()
    srv := server.NewServer(db)

    srv.RegisterModules()
    srv.Run()
}
```

---

## 🌱 Server with Auto Module Loading (`internal/server/server.go`)

```go
package server

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "your-project/internal/user"
)

type Server struct {
    Engine *gin.Engine
    DB     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
    return &Server{
        Engine: gin.Default(),
        DB:     db,
    }
}

func (s *Server) RegisterModules() {
    api := s.Engine.Group("/api")
    userModule := user.NewModule(s.DB)
    userModule.RegisterRoutes(api)
}

func (s *Server) Run() {
    s.Engine.Run(":8080")
}
```

---

## 🔧 Config File (`config/config.go`)

```go
package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment")
    }
}
```

---

## 📦 How This Helps

✔ Highly modular – easy to split into microservices later  
✔ Clear separation of concerns  
✔ No module needs to know internal structure of others  
✔ Supports DI-like architecture and single instances of services/repositories  

---

## ✅ Next Steps (Opcional)

- Add test files (`*_test.go`)  
- Add Swagger/OpenAPI documentation  
- Add Dockerfile & docker-compose  
- Add authentication (JWT)  
- Add Makefile  

---

Let me know if you want those generated too!
