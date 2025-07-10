# AGENTS.md - Coding Guidelines for Core API

## Build/Test Commands
- Build: `task build` (rebuilds Docker container)
- Run: `task up` (starts service and dependencies)
- Migrate: `task migrate` (applies database migrations)
- Seed: `task seed` (seeds database with test data)
- Init: `task init` (full setup: up, migrate, seed)
- Logs: `task logs` (view all service logs)
- Lint: `golangci-lint run` (uses .golanci.yml config with all linters enabled)
- No test files exist yet - create with `*_test.go` suffix following Go conventions

## Code Style & Conventions
- **Imports**: Standard library first, then third-party, then local packages
- **Naming**: Use camelCase for functions/variables, PascalCase for exported types
- **Error Handling**: Always wrap errors with context using `fmt.Errorf("message: %w", err)`
- **Types**: Use interfaces for service contracts, structs for implementations
- **JSON Tags**: Use snake_case for JSON field names in request/response structs
- **Validation**: Use struct tags for Huma validation (minLength, maxLength, format)
- **Context**: Always pass `context.Context` as first parameter to functions
- **UUIDs**: Use `github.com/google/uuid` for ID generation and handling
- **Database**: Use `jmoiron/sqlx` for database operations, `jackc/pgx/v5` for PostgreSQL
- **API Framework**: Using Huma v2 with Echo v4 for HTTP routing and validation-model