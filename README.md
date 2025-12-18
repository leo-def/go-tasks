# GoTasks


## Table of Contents
- [Commands](#commands)
- [Makefile Commands](#makefile-commands)
- [Debugging with Hot Reload](#debugging-with-hot-reload)
- [Logging](#logging)
- [Environment Files](#environment-files)
- [JWT Secret Generator](#jwt-secret-generator)
- [Project Structure](#project-structure)
- [Docker](#docker)
- [Migrations](#migrations)
- [API Docs](#api-docs)
- [Testing](#testing)
- [Schemas](#schemas)
- [Endpoints](#endpoints)
- [Rules](#rules)
- [Lifecycle Flow](#lifecycle-flow)
- [Test Summary](#test-summary)


## Commands
- GO run: `go run ./cmd/api/main.go`
- GO build: `go build -o go-tasks ./cmd/api/main.go`
- GO Install: `go install ./cmd/api/main.go`
- Init Doc: `swag init -g ./cmd/api/main.go -o ./docs`

## Makefile Commands
- Run: `make run`
- Dev (hot reload): `make dev`
- **Debug (hot reload + debugging): `make debug`**
- Debug (headless, no hot reload): `make debug-headless`
- Debug Stop: `make debug-stop`
- Build: `make build`
- Install: `make install`
- Docs: `make docs`
- Test: `make test`
- Test (race): `make test-race`
- Test (coverage): `make test-cover`
- Test (HTML coverage): `make test-cover-html`
- Help: `make help`
- Migrate Create: `make migrate-create name=add_table`
- Migrate Up: `DATABASE_URL=postgres://user:pass@host:port/dbname make migrate-up`
- Migrate Down: `DATABASE_URL=... make migrate-down`
- Migrate Status: `DATABASE_URL=... make migrate-status`
- Migrate Reset: `DATABASE_URL=... make migrate-reset`

## Debugging with Hot Reload

This project supports debugging with hot reload using Air and Delve. You can set breakpoints in your IDE and have the application automatically restart when you make code changes.

### Quick Start
```bash
# Start the application with hot reload and debugging
make debug
```

The application will:
- ✅ Run with hot reload (detects file changes and restarts)
- ✅ Start with Delve debugger on port `:2345`
- ✅ Serve the API on port `:8080`
- ✅ Automatically continue execution (no manual intervention needed)

### IDE Setup (VS Code)

1. **Use the provided launch configuration**: A `.vscode/launch.json` file is already configured
2. **Set breakpoints**: Click in the gutter next to line numbers in your Go files
3. **Connect debugger**: 
   - Go to Run and Debug panel (Ctrl+Shift+D)
   - Select "Debug Go App" configuration
   - Click "Start Debugging" or press F5
   - Alternatively, when `make debug` is running, select "Attach to Delve (Air debug)" to attach to the headless server on `localhost:2345`

### Manual Debugger Connection

If you prefer to connect manually:
```bash
# The debugger is available at localhost:2345
# You can connect using any DAP-compatible debugger client
```

### Testing Breakpoints

1. Start debugging: `make debug`
2. Set a breakpoint in `internal/auth/middleware.go` at line 15 or 45
3. Trigger the middleware with a request:
   ```bash
   curl -H "Authorization: Bearer test-token" http://localhost:8080/api/v1/admin/company/
   ```
4. Your breakpoint should be hit, and you can inspect variables

### Hot Reload in Action

1. Start debugging: `make debug`
2. Make any change to a `.go` file
3. Save the file
4. Watch the terminal - Air will detect the change, rebuild, and restart with debugging enabled
5. Your breakpoints remain active after restart

### Available Debug Commands

- `make debug` - Hot reload with debugging (recommended for development)
- `make debug-headless` - Debug without hot reload (for production debugging)
- `make debug-stop` - Stop all debug processes

### Troubleshooting

**Port already in use?**
```bash
make debug-stop  # Stop any existing debug processes
make debug       # Start fresh
```

**Breakpoints not working?**
- Ensure your IDE is connected to `localhost:2345`
- Check that the `.vscode/launch.json` configuration is correct
- Verify the application is running with `curl http://localhost:8080/health`

## Logging

- Structured JSON logs with levels controlled by `LOG_LEVEL` (`debug`, `info`, `warn`, `error`). Default is `info`.
- Minimal by default: only important events (startup, warnings, errors) are logged.
- Optional HTTP request logs: set `HTTP_LOG=true` to enable Gin's request logger.
- Example `.env`:
  ```
  LOG_LEVEL=info
  HTTP_LOG=false
  ```
 
The logger lives in `internal/pkg/logger/`. Use `logger.Info/Warn/Error` and `logger.WithContext(c, ...)` to emit tracing-friendly logs.

## Cascade Deletion

- Company -> Activity: deleting a `company` removes its `activities` (DB cascade)
- Company -> Collaborator: deleting a `company` removes its `collaborators` (DB cascade)
- Activity -> Task: deleting an `activity` removes its `tasks` (DB cascade)
- Activity -> Participation: deleting an `activity` removes its `participations` (DB cascade)
- Participation -> Rating: deleting a `participation` removes its `ratings` (DB cascade)
- Task -> Assignment: deleting a `task` removes its `assignments` (DB cascade)
- Collaborator -> Participation: deleting a `collaborator` removes their `participations` (DB cascade)
- Collaborator -> Assignment: assignments where the `collaborator` is the assigner are removed before deleting the collaborator (service logic)
- Collaborator -> Account: if the deleted collaborator’s `account` has no other collaborator links and the account role is `USER`, the account is deleted (service logic)

This behavior is recursive. For example, deleting an `activity` also removes related `tasks` and `participations`, which removes dependent `assignments` and `ratings` automatically.

## Environment Files

- Use `.env.example` as the canonical base; copy it to `.env` and set `JWT_SECRET` to a strong random value.
- For specific environments, copy the base and apply overrides:
  - `.env.postgres.example`: sets `DATABASE_DRIVER=postgres`, `DATABASE_URL`, and `HTTP_LOG=true`.
  - `.env.sqlite3.example`: notes that base defaults already match SQLite.
  - `.env.docker.example`: sets `APP_ENV=dev`, Postgres DSN for container, expanded `TRUSTED_PROXIES`, and `HTTP_LOG=true`.
- This approach removes duplicated variables across example files and keeps configuration clear.

## JWT Secret Generator

- Generate a strong secret from a Go CLI:
  - `make secret`
  - Options:
    - `-bytes=64` to increase entropy
    - `-format=hex|base64|base64url` to change output
    - `-env` to print as `JWT_SECRET=...` line
- Paste into your `.env` or `.env.docker`.
- Prefer different secrets per environment.

**Hot reload not detecting changes?**
- Check that your file is not in the excluded directories (tmp, vendor, testdata, docs, migrations, db)
- Ensure you're editing `.go` files (other extensions may not trigger reload)

## Project Structure
See [STRUCTURE.md](./STRUCTURE.md) for the project layout and package descriptions.


## Docker

### Run without hot reload (production-like)
- Prerequisites:
  - Copy env file: `cp .env.docker.example .env.docker` and edit values
- Command:
  - `docker compose up --build`
- What it does:
  - Builds the app using the multi-stage `Dockerfile`
  - Loads environment from `.env.docker`
  - Runs compiled binary (`go-tasks`) inside Alpine

### Run with hot reload (development)
- Prerequisites:
  - Copy env file: `cp .env.docker.example .env.docker` and edit values
- Command:
  - `docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build`
- What it does:
  - Builds using the `dev` stage (Go toolchain installed)
  - Mounts your source code into the container
  - Installs Air and runs with hot reloading
  - Watches `.go` files and restarts on changes

### Notes
- Environment is loaded from `.env.docker` (do not commit real secrets).
- For local, non-Docker dev, use `make dev` to run Air on your host.
- Migrations inside Docker dev are applied via `make migrate-up` in the container.


## Migrations
- Tooling: uses `goose` via `go run github.com/pressly/goose/v3/cmd/goose@latest`.
- Files: SQL migrations live in `./migrations` as timestamped files.
- Env: set `DATABASE_URL` to your Postgres DSN (same as app uses).

### Examples
- Create a migration: `make migrate-create name=add_column_to_companies`
- Apply up: `DATABASE_URL=postgres://user:pass@localhost:5432/mydb make migrate-up`
- Rollback one: `DATABASE_URL=... make migrate-down`
- Status: `DATABASE_URL=... make migrate-status`
- Reset all: `DATABASE_URL=... make migrate-reset`

## Schemas 

- enum Role
```
    ADMIN
    OPS
    USER
```

- enum CollaboratorRole
```
    OWNER
    MANAGER
    OPS
```

- enum Status
```
    PENDING
    INITIALIZED
    DELIVERED
    CANCELLED
```

- model Company
```
    ID uuid.ID
    Title string
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Account
```
    ID uuid.ID
    name string
    username string
    password string
    email string
    phone string
    role: Role
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Collaborator
```
    ID uuid.ID
    Title string
    role: CollaboratorRole
    accountId: uuid
    companyId: uuid
    CreatedAt time.Time
    updatedAt: DateTime
```
  
- model StatusUpdate 
```
    ID uuid.ID
    lifecycleId: uuid
    statusBefore: Status
    statusAfter: Status
    updateDate: Date
    CreatedAt time.Time
    updatedAt: DateTime
```
 
- model Lifecycle
```
    ID uuid.ID
    initDate: Date
    dueDate: Date
    status: Status
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Activity
```
    ID uuid.ID
    Title string
    lifecycleId: uuid
    companyId: uuid
    ownerId: uuid
    createdBy: uuid
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Task
```
    ID uuid.ID
    Title string
    activityId: uuid
    lifecycleId: uuid
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Participation
```
    ID uuid.ID    
    collaboratorId: uuid
    activityId: uuid
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Assignment
```
    ID uuid.ID
    participationId: uuid
    assigner: uuid
    taskId: uuid
    assignDate: Date
    CreatedAt time.Time
    updatedAt: DateTime
```

- model Rating
```
    ID uuid.ID
    rate: float
    participationId: uuid
    collaboratorId: uuid
    CreatedAt time.Time
    updatedAt: DateTime
```

## API Docs

- OpenAPI specs live under `./docs`:
  - `./docs/swagger.yaml`
  - `./docs/swagger.json`
- Regenerate docs: `make docs` or `swag init -g ./cmd/api/main.go -o ./docs`
- When the API is running locally (`make run` or `make dev`), the spec reflects current routes and models.

## Testing

- Quick run all packages:
  - `make test`
- With race detector:
  - `make test-race`
- With coverage summary:
  - `make test-cover`
- Generate HTML coverage report:
  - `make test-cover-html`
- Direct `go test` examples:
  - Run all: `go test ./...`
  - Verbose in a package: `go test ./internal/pkg/httpx -v`
  - Coverage profile: `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`

## Run Locally (SQLite3)

- Prereqs
  - Install tools: `go install github.com/swaggo/swag/cmd/swag@latest`, `go install github.com/pressly/goose/v3/cmd/goose@latest`
  - Optional hot reload: `go install github.com/air-verse/air@latest`
- Env
  - Copy local env: `cp .env.local.example .env`
  - Ensure `DATABASE_DRIVER=sqlite3` and `DATABASE_URL=./db/go-tasks.db`
  - Create DB dir: `mkdir -p db`
- Generate docs
  - `make docs`
- Migrate (SQLite3)
  - `make migrate-up-sqlite`
- Run
  - Hot reload: `make dev`
  - Basic run: `make run`
- Verify
  - Swagger UI: `http://localhost:8080/swagger/index.html`
  - Health: `http://localhost:8080/health`

## Run with Docker Compose (Postgres)

- Env
  - Copy docker env: `cp .env.docker.example .env.docker`
  - `.env.docker` uses `DATABASE_DRIVER=postgres` and a Postgres `DATABASE_URL`
- Start dev stack
  - `docker compose -f docker-compose.dev.yml up --build`
  - The dev container runs: `make docs`, `make migrate-up` (Postgres), and starts with hot reload via Air
- Verify
  - Swagger UI: `http://localhost:8080/swagger/index.html`
  - Health: `http://localhost:8080/health`


## Endpoints
Permissions and validations overview:
- Company operation: JWT must include `company_id`. If `systemRole=USER`, the request must also satisfy `companyRole` checks.
- Lifecycle validation: child lifecycles must start on or after the parent start date and finish on or before the parent due date.

Legend:
- Access levels:
  - System/Public: no JWT; no role checks.
  - System/Protected: JWT required; `systemRole` enforced where applicable.
  - Company Operation: JWT with `company_id` required; if `systemRole=USER`, `companyRole` is enforced by endpoint.
- Roles:
  - `systemRole`: `ADMIN`, `OPS`, `USER`.
  - `companyRole`: `OWNER`, `MANAGER`, `OPS`.
- Company context:
  - `company_id` must be present in JWT for Company Operation endpoints.

### Auth
Permissions: `systemRole=public` (no auth) and `systemRole=protected` (JWT required)

#### Auth.Public:
Permissions:
- Access: Public

| Method | Path             | Operation  | System Role | Company Role | Notes | Task   |
|--------|------------------|------------|-------------|--------------|-------|--------|
| POST   | `/auth/sign-in`  | SignIn     | PUBLIC      |              |       |  IMPL  |
| POST   | `/auth/sign-up`  | SignUp     | PUBLIC      |              |       |  IMPL  |

#### Auth.Protected:
Permissions:
- Access: Protected
- CompanyContext: not required
- systemRole: `admin|ops|user`
- companyRole: n/a

| Method | Path              | Operation       | System Role    | Company Role | Notes | Task   |
|--------|-------------------|-----------------|----------------|--------------|-------|--------|
| GET    | `/auth/me`        | GetMe           | ADMIN-OPS-USER |              |       |  IMPL  |
| PUT    | `/auth/password`  | UpdatePassword  | ADMIN-OPS-USER |              |       |  IMPL  |
| PUT    | `/auth/email`     | UpdateEmail     | ADMIN-OPS-USER |              |       |  IMPL  |
| PUT    | `/auth/phone`     | UpdatePhone     | ADMIN-OPS-USER |              |       |  IMPL  |

### Company
Permissions:
- Access: Protected
- CompanyContext: not required
- systemRole: `admin|ops`
- companyRole: n/a

| Method | Path            | Operation | System Role | Company Role | Notes | Task   |
|--------|-----------------|-----------|-------------|--------------|-------|--------|
| GET    | `/company/`     | GetAll    | ADMIN-OPS   |              |       |  IMPL  |
| GET    | `/company/:id`  | GetById   | ADMIN-OPS   |              |       |  IMPL  |
| POST   | `/company/`     | Create    | ADMIN-OPS   |              |       |  IMPL  |
| PUT    | `/company/:id`  | Update    | ADMIN-OPS   |              |       |  IMPL  |
| DELETE | `/company/:id`  | Delete    | ADMIN-OPS   |              |       |  IMPL  |

### Collaborator

#### Collabrator Owner
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: `owner`

| Method | Path                    | Operation          | System Role | Company Role | Notes | Task   |
|--------|-------------------------|--------------------|-------------|--------------|-------|--------|
| GET    | `/collaborator/`        | GetAll             |             | OWNER        |       |  IMPL  |
| GET    | `/collaborator/:id`     | GetById            |             | OWNER        |       |  IMPL  |
| POST   | `/collaborator/`        | Create             |             | OWNER        |       |  IMPL  |
| POST   | `/collaborator/account` | Create With Account|             | OWNER        |       |  IMPL  |
| PUT    | `/collaborator/:id`     | Update             |             | OWNER        |       |  IMPL  |
| DELETE | `/collaborator/:id`     | Delete             |             | OWNER        |       |  IMPL  | 

### Collaborator Protected
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: any

| Method | Path                      | Operation | System Role    | Company Role | Notes | Task   |
|--------|---------------------------|-----------|----------------|--------------|-------|--------|
| GET    | `/collaborator/my-roles`  | MyRoles   | ADMIN-OPS-USER |              |       |  IMPL  |

### Collaborator SubCollaborator
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: `manager|owner`

| Method | Path                     | Operation | System Role | Company Role  | Notes | Task   |
|--------|--------------------------|-----------|-------------|---------------|-------|--------|
| GET    | `/collaborator/sub/`     | GetAll    |             | MANAGER-OWNER |       |  IMPL  |
| GET    | `/collaborator/sub/:id`  | GetById   |             | MANAGER-OWNER |       |  IMPL  |
| POST   | `/collaborator/sub/`     | Create    |             | MANAGER-OWNER |       |  IMPL  |
| PUT    | `/collaborator/sub/:id`  | Update    |             | MANAGER-OWNER |       |  IMPL  |
| DELETE | `/collaborator/sub/:id`  | Delete    |             | MANAGER-OWNER |       |  IMPL  |


### Participation 
Permissions:

- Access: Company Operation (self)
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: `manager|owner|ops`

| Method | Path                 | Operation                 | System Role | Company Role  | Notes | Task   |
|--------|----------------------|---------------------------|-------------|---------------|-------|--------|
| GET    | `/participation/my`  | Get User Participations   |             | MANAGER-OWNER |       |  IMPL  |

### Activity
Permissions:
- Access: Company Operation
- CompanyContext: `collaborator_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole : `manager|owner`

| Method | Path             | Operation | System Role | Company Role  | Notes | Task   |
|--------|------------------|-----------|-------------|---------------|-------|--------|
| GET    | `/activity/`     | GetAll    |             | MANAGER-OWNER |       |  IMPL  |
| GET    | `/activity/:id`  | GetById   |             | MANAGER-OWNER |       |  IMPL  |
| POST   | `/activity/`     | Create    |             | MANAGER-OWNER |       |  IMPL  |
| PUT    | `/activity/:id`  | Update    |             | MANAGER-OWNER |       |  IMPL  |
| DELETE | `/activity/:id`  | Delete    |             | MANAGER-OWNER |       |  IMPL  |


### Activity Status
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole (for USER): `manager|owner`

| Method | Path                             | Operation                | System Role | Company Role  | Notes                                 | Task   |
|--------|----------------------------------|--------------------------|-------------|---------------|---------------------------------------|--------|
| PUT    | `/activity/:activity_id/status/` | Update Activity Status   |             | MANAGER-OWNER | Validates child lifecycle constraints |  IMPL  |

### Activity Participation
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: `manager|owner`

| Method | Path                                    | Operation             | System Role | Company Role  | Notes | Task   |
|--------|-----------------------------------------|-----------------------|-------------|---------------|-------|--------|
| GET    | `/participation/activity/:activity_id/` | GetAll Participation  |             | MANAGER-OWNER |       |  IMPL  |
| POST   | `/participation/activity/:activity_id/` | Create Participation  |             | MANAGER-OWNER |       |  IMPL  | 



### Activity Participation Rating
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole: `manager|owner`

| Method | Path                                        | Operation         | System Role | Company Role  | Notes | Task   |
|--------|---------------------------------------------|-------------------|-------------|---------------|-------|--------|
| POST   | `/rating/participation/:participation_id/`  | Create Rating     |             | MANAGER-OWNER |       |  IMPL  |  

### Activity Task
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole (for USER): `manager|owner`

| Method | Path                                | Operation | System Role | Company Role  | Notes | Task   |
|--------|-------------------------------------|-----------|-------------|---------------|-------|--------|
| GET    | `/activity/:activity_id/task/`      | GetAll    |             | MANAGER-OWNER |       |  IMPL  |
| GET    | `/activity/:activity_id/task/:id`   | GetById   |             | MANAGER-OWNER |       |  IMPL  |
| POST   | `/activity/:activity_id/task/`      | Create    |             | MANAGER-OWNER |       |  IMPL  |
| PUT    | `/activity/:activity_id/task/:id`   | Update    |             | MANAGER-OWNER |       |  IMPL  |
| DELETE | `/activity/:activity_id/task/:id`   | Delete    |             | MANAGER-OWNER |       |  IMPL  |
| PUT    | `/activity/:activity_id/task/:id/status`   | Update Status    |             | MANAGER-OWNER |       |  IMPL  |

My Tasks
Permissions:
- Access: Company Operation (self)
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole (for USER): `manager|owner|ops`

| Method | Path                              | Operation     | System Role | Company Role      | Notes | Task   |
|--------|-----------------------------------|---------------|-------------|-------------------|-------|--------|
| GET    | `/activity/my/:activity_id/task/` | Get My Tasks  |             | MANAGER-OWNER-OPS | Only if the activity includes my participation; returns only my tasks | IMPL |

### Task Assignment
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole (for USER): `manager|owner`

| Method | Path                                        | Operation    | System Role | Company Role  | Notes | Task   |
|--------|---------------------------------------------|--------------|-------------|---------------|-------|--------|
| POST   | `/task/:task_id/assignment/`               | Create       |             | MANAGER-OWNER |       |  IMPL  |
| GET    | `/task/:task_id/assignment/`               | FindByTaskId |             | MANAGER-OWNER |       |  IMPL  |  
| GET    | `/task/:task_id/assignment/:id`            | FindById     |             | MANAGER-OWNER |       |  IMPL  |
| PATCH  | `/task/:task_id/assignment/:id/deactivate` | Deactivate   |             | MANAGER-OWNER |       |  IMPL  |
| DELETE | `/task/:task_id/assignment/:id`            | Delete       |             | MANAGER-OWNER |       |  IMPL  |

### Task Status
Permissions:
- Access: Company Operation
- CompanyContext: `company_id` required
- systemRole: `admin|ops`; `user`: must satisfy companyRole
- companyRole (for USER): `manager|owner`

| Method | Path                                        | Operation               | System Role | Company Role  | Notes | Task   |
|--------|---------------------------------------------|-------------------------|-------------|---------------|-------|--------|
| PUT    | `/task/:id/status/`   | Update Task Status      |             | MANAGER-OWNER | Validates parent lifecycle constraints |       |  IMPL  |


## Lifecycle Flow
- CREATED → BLOCKED | IN_PROGRESS | CANCELLED
- IN_PROGRESS → BLOCKED | COMPLETED | CANCELLED
- BLOCKED → CANCELLED | IN_PROGRESS

## Rules

**Activity**
- Create
  - Status is `CREATED`
  - `dueDate` must be after `initDate`
- Edit
  - Status cannot be edited here
  - `dueDate` must be after `initDate`
  - `startDate` must be before all tasks `startDate`
  - `endDate` must be after all tasks `endDate`
- Delete
  - Activity must be in an inactive status (`COMPLETED`, `CANCELLED`)
  - All tasks must be in an inactive status (`COMPLETED`, `CANCELLED`)
- Status Update
  - Status must follow lifecycle flow
  - If updating to an inactive status, all tasks must already be inactive (`COMPLETED`, `CANCELLED`)

**Task**
- Create
  - Status is `CREATED`
  - `dueDate` must be after `initDate`
  - Activity status must be active (`CREATED`, `IN_PROGRESS`, `BLOCKED`)
  - `startDate` must be between activity `startDate` and `dueDate`
  - `dueDate` must be between activity `startDate` and `endDate`
- Edit
  - Status cannot be edited here
  - `dueDate` must be after `initDate`
  - `startDate` must be between activity `startDate` and `dueDate`
  - `dueDate` must be between activity `startDate` and `endDate`
- Delete
  - Task must be in an inactive status (`COMPLETED`, `CANCELLED`)
- Status Update
  - Status must follow lifecycle flow
  - If updating to an active status, activity must be active (`CREATED`, `IN_PROGRESS`, `BLOCKED`)

## Test Summary
- Admin / Company: Complete
- Owner / Collaborator: Complete
- Collaborator / Subordinates: Complete
  - Actions can be performed only over subordinates: Complete
  - Only subordinates can be fetched: Complete
- Activity: Pending
  - Activity with pending tasks cannot be deleted: Pending
- Activity / Status: Partial
  - Status update must follow status flow: Complete
  - Task status must be considered: Pending
- Activity / Task: Pending
  - Activity status must be considered: Pending
  - Activity must be active: Pending
- Activity / Task / Status: Partial
  - Status update must follow status flow: Complete
  - Activity status must be considered: Pending
- Activity / Task / Assignment: Pending
  - Task and Activity must be active: Pending
  - Create participation for activity if it does not exist: Pending
- Collaborator / Rating: Pending
  - Must have a participation for activity: Pending
