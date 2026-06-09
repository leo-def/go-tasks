# GoTasks - Technical Specification

> Go 1.25 + Gin + GORM REST API for company task and activity management with JWT authentication.
> Demonstrates modular Go architecture with 12 isolated domain modules.

## Executive Summary

GoTasks is a **Go 1.25 + Gin + GORM** REST API for managing company operations: accounts, collaborators, companies, activities, tasks, assignments, lifecycles, ratings, and participations. It uses a modular architecture where each domain is a self-contained package (`model.go`, `repository.go`, `service.go`, `handler/controller.go`, `module.go`). Authentication is JWT-based. Database is PostgreSQL in production, SQLite in testing. API documentation is auto-generated via Swagger.

---

## 1. Problem Statement

### Context
A collaborative task management platform for companies: companies onboard staff as collaborators, create activities with lifecycle tracking, assign tasks to collaborators, and rate performance.

### Goals
- Company and collaborator management with role-based access
- Activity and task lifecycle tracking (status: pending → in_progress → done)
- JWT authentication with role-based permissions (Admin, Ops, Owner, Manager)
- Task assignment and performance rating
- OpenAPI/Swagger documentation
- PostgreSQL production DB + SQLite for tests

### Success Metrics
- [x] 12 isolated domain modules
- [x] JWT-based RBAC (Admin/Ops/Owner/Manager roles)
- [x] Database transactions for atomic operations
- [x] Swagger/OpenAPI auto-generated docs
- [x] PostgreSQL + SQLite support (GORM dialects)
- [ ] 85%+ test coverage on service layer
- [ ] Response time p99 < 50ms

---

## 2. Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.25.2 |
| HTTP Router | Gin | v1.11.0 |
| ORM | GORM | v1.31.1 |
| DB (production) | PostgreSQL | 15+ |
| DB (testing) | SQLite | v1.6.0 |
| Auth | golang-jwt/jwt | v5.3.0 |
| IDs | Google UUID | v1.6.0 |
| API Docs | Swaggo (Swagger) | v1.16.6 |
| Hashing | golang.org/x/crypto | v0.45.0 |
| Hot Reload | Air | Latest |
| Debugger | Delve | Latest |
| Testing | Testify | v1.11.1 |
| Mocking | go.uber.org/mock | v0.5.0 |

---

## 3. Architecture

```
┌────────────────────────────────────────────┐
│         HTTP Clients                        │
└─────────────────┬──────────────────────────┘
                  │ HTTP/REST
                  ▼
┌────────────────────────────────────────────┐
│   Gin Router + Middleware                   │
│   (JWT Auth, Logging, Error Handling, CORS) │
└─────────────────┬──────────────────────────┘
                  │
       ┌──────────┴──────────┐
       ▼                     ▼
 Module Layer          Shared Layer (pkg/)
 (12 domains)          logger, response, httpx
       │
       ▼
 Service + Repository Layer
 (Business logic + GORM queries)
       │
       ▼
 PostgreSQL (GORM, UUID primary keys)
```

---

## 4. Module Registry (12 Domains)

| Module | Purpose | Key Files |
|--------|---------|-----------|
| `account` | User accounts (login, register) | model, repository, service, controller |
| `auth` | JWT authentication + middleware | JWT generation, middleware |
| `collaborator` | Company staff with roles | model (Role, AccountID, CompanyID) |
| `company` | Company management | CRUD + listing |
| `activity` | Activities within a company | Activity + lifecycle + owner/collaborator controllers |
| `task` | Tasks within an activity | Create/update/delete with lifecycle |
| `lifecycle` | Lifecycle state machine | InitDate, DueDate, Status, StatusUpdates |
| `assignment` | Task assignment to collaborators | AssignTask, UnassignTask |
| `participation` | User participation in activities | Join/leave activities |
| `rating` | Performance ratings | Rate collaborator performance |
| `health` | Health check endpoint | `GET /health` |
| `app` | App orchestration + seeding | Routes registration, DB seed |

---

## 5. Data Models (Key Structures)

```go
// TaskInfo — stored in "tasks" table
type TaskInfo struct {
    ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Title       string     `gorm:"not null"`
    ActivityID  uuid.UUID  `gorm:"type:uuid;not null"`
    LifecycleID uuid.UUID  `gorm:"type:uuid;not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   time.Time  // ⚠️ Bug: should be gorm.DeletedAt
}

// TaskLifecycle — "lifecycles" table
type TaskLifecycle struct {
    ID       uuid.UUID
    InitDate time.Time    `gorm:"not null"`
    DueDate  time.Time    `gorm:"not null"`
    Status   string       `gorm:"not null"`
    Updates  []TaskStatusUpdate
}

// TaskCollaborator — "collaborators" table
type TaskCollaborator struct {
    ID        uuid.UUID
    Role      string
    AccountID uuid.UUID
    CompanyID uuid.UUID
    Account   TaskAccount  // ⚠️ Bug: exposes password hash in JSON
}
```

---

## 6. API Endpoints

```
# Auth
POST /api/v1/auth/login
POST /api/v1/auth/register

# Account
GET/POST/PUT/DELETE /api/v1/accounts
GET /api/v1/accounts/:id

# Company
GET/POST/PUT/DELETE /api/v1/companies
GET /api/v1/companies/:id

# Collaborator
GET/POST/PUT/DELETE /api/v1/collaborators
GET /api/v1/collaborators/:id

# Activity
GET/POST/PUT/DELETE /api/v1/activities
GET /api/v1/activities/:id
GET /api/v1/companies/:id/activities   ← by company

# Task
GET/POST/PUT/DELETE /api/v1/tasks
GET /api/v1/activities/:id/tasks       ← by activity
PUT /api/v1/tasks/:id/status           ← status update

# Assignment
POST/DELETE /api/v1/assignments

# Rating
GET/POST /api/v1/ratings

# Health
GET /api/v1/health

# Swagger
GET /swagger/index.html
```

---

## 7. Task Creation (Transactional Pattern)

```go
// service-activity.go — atomic task creation
func (s *serviceActivity) Create(task *Task, activityID uuid.UUID) error {
    tx := s.repository.GetConnection().BeginTransaction()
    // 1. Create lifecycle record
    if err := s.lifecycleService.CreateTx(tx, &lifecycle); err != nil {
        tx.Rollback(); return err
    }
    // 2. Create task with lifecycle FK
    if err := s.repository.CreateInfoTx(tx, &task.TaskInfo); err != nil {
        tx.Rollback(); return err
    }
    return tx.Commit()
}
```

---

## 8. Testing Strategy

```bash
make test              # All tests
make test-race         # Race condition detection
make test-cover        # Coverage report
make test-cover-html   # HTML coverage report
```

- Test database: SQLite (GORM `gorm.io/driver/sqlite`)
- Mocking: `go.uber.org/mock`
- Assertions: `testify/assert`

---

## 9. Deployment & Operations

```bash
make run        # Build and run
make dev        # Hot reload with Air
make debug      # Delve debugger

# Docker
docker-compose up
```

**Env vars:** `DATABASE_URL`, `JWT_SECRET`, `PORT`, `GIN_MODE`, `LOG_LEVEL`

---

## 10. Issues Found

### Critical Bugs
- **`DeletedAt time.Time` (not `gorm.DeletedAt`)** in `TaskInfo`, `TaskActivity`, `TaskLifecycle`, and `TaskCollaborator` models — GORM's soft-delete feature requires the field type to be `gorm.DeletedAt` (a pointer-based nullable type with `Valid bool`). Using plain `time.Time` means GORM's `Unscoped` / soft-delete queries will not work as expected; deletes may permanently remove records or soft-delete won't filter properly.

- **`TaskCollaborator` embeds `TaskAccount` which contains password hash** — if `TaskAccount` has a `Password` or `PasswordHash` field with `json:"-"` missing, the password hash is serialized in API responses containing collaborator data.

- **Foreign key typo in `TaskActivity.Lifecycle`**: `gorm:"foreignKey:LifecycleID;reference:ID"` — should be `references:ID` (plural). The singular form is silently ignored by GORM, potentially causing incorrect join behavior.

### Security
- Password hash should have `json:"-"` tag on the account model to prevent serialization.
- JWT secret should be validated at startup — if `JWT_SECRET` is empty, tokens would be signed with an empty key.

### Performance
- No pagination default limit — list endpoints could return unbounded result sets.
