# Core Service

This is the core backend service that MatchAlly relies on. It provides essential functionality for the platform.

## Prerequisites

- Go 1.24 or later
- Docker

## Getting Started

### 1. Setup Dependencies

The service requires PostgreSQL and Redis. You can start these services using Docker Compose:

```bash
docker compose up -d
```

This will start:
- PostgreSQL 17 on port 5432
- Redis (Valkey) on port 6379

### 3. Database Migrations

The project uses Goose for database migrations. To run migrations:

```bash
# Run all pending migrations
goose -dir migrations postgres "postgres://core:secret@localhost:5432/core?sslmode=disable" up

# To rollback the last migration
goose -dir migrations postgres "postgres://core:secret@localhost:5432/core?sslmode=disable" down
```

### 4. Running the Service

To start the service in development mode:

```bash
go run main.go serve
```

The service will be available at `http://localhost:8080` by default.

## Development

### Project Structure

- `cmd/` - Command-line interface and service entry points
- `internal/` - Private application code
- `migrations/` - Database migration files
- `seeds/` - Database seed data
- `docker-compose.yaml` - Local development environment configuration
- `taskfile.yaml` - Common development tasks

### Development Tools

#### Goose
[Goose](https://github.com/pressly/goose) is used for database migrations. Our migrations are written to be compatible with this tool.

#### go-mod-upgrade
[go-mod-upgrade](https://github.com/oligot/go-mod-upgrade) helps manage Go dependencies:
```bash
# Check for dependency updates
go-mod-upgrade

# After updating dependencies, run:
go mod tidy
```
