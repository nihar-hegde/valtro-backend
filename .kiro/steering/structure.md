# Project Structure & Architecture

## Directory Organization

### Application Entry Point
- `cmd/api/main.go` - Main application entry point, handles initialization

### Internal Package Structure
```
internal/
├── database/           # Database connection and configuration
├── handlers/          # HTTP request handlers (controllers)
│   ├── health/       # Health check endpoints
│   ├── organization/ # Organization CRUD operations
│   ├── project/      # Project CRUD operations
│   └── user/         # User CRUD operations
├── middleware/        # HTTP middleware (currently empty)
├── models/           # Data models and database schemas
├── repositories/     # Data access layer (currently empty)
├── server/           # Server setup and routing
└── services/         # Business logic layer
    ├── organization/
    ├── project/
    └── user/
```

## Architecture Patterns

### Handler Pattern
- Each domain (user, organization, project) has its own handler package
- Handlers are structs that hold database dependencies
- Constructor pattern: `NewHandler(db *gorm.DB) *Handler`
- Standard CRUD methods: Create, GetAll, GetByID, Update, Delete

### Server Structure
- `Server` struct holds all dependencies (database, handlers)
- Centralized route registration in `RegisterRoutes()`
- Chi router with middleware (Logger, Recoverer)

### Model Conventions
- UUID primary keys with PostgreSQL `gen_random_uuid()`
- GORM struct tags for database constraints
- Soft deletion support with `gorm.DeletedAt`
- Comprehensive field documentation with comments
- Nullable fields use pointer types (`*string`, `*time.Time`)

### API Structure
- Versioned API: `/api/v1/`
- RESTful resource routing
- Consistent response format with JSON
- Standard HTTP methods and status codes

## Naming Conventions
- Package names: lowercase, singular
- Handler methods: PascalCase matching HTTP operations
- Database fields: PascalCase with GORM tags
- Environment variables: UPPERCASE_SNAKE_CASE