# Technology Stack

## Core Technologies
- **Language**: Go 1.24.4
- **Web Framework**: Chi v5 (lightweight HTTP router)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: Clerk integration
- **Environment Management**: godotenv for .env file loading

## Key Dependencies
- `github.com/go-chi/chi/v5` - HTTP router and middleware
- `gorm.io/gorm` - ORM for database operations
- `gorm.io/driver/postgres` - PostgreSQL driver for GORM
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading

## Build System & Commands

### Make Commands
```bash
make build    # Build the application binary
make run      # Run the application directly
make clean    # Clean built binaries and temp files
make watch    # Live reload development with Air
```

### Development Workflow
- Use `make watch` for development with live reload
- Air configuration in `air.toml` watches `.go` files and rebuilds automatically
- Main entry point: `cmd/api/main.go`
- Binary output: `./main` (production) or `./tmp/main` (development)

### Environment Setup
- Requires `DATABASE_URL` environment variable
- Optional `PORT` environment variable (defaults to 8080)
- Uses `.env` file for local development (loaded via godotenv)