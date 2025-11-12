# Backend Order System (BO)

A modular Go backend API using PostgreSQL, Chi router, and sqlc for type-safe database queries.

## Features

- 🏗️ **Modular Architecture**: Independent modules for users and orders
- 🔒 **Type-Safe Queries**: Using sqlc for generated database code
- 🚀 **Fast HTTP Router**: Chi router with middleware support
- 🐘 **PostgreSQL**: Database with Docker support
- 📝 **Structured Logging**: JSON logging with slog
- ✅ **Validation**: Request validation with go-playground/validator
- 🔄 **Graceful Shutdown**: Proper server lifecycle management
- 📚 **Swagger Documentation**: Interactive API documentation with Swagger UI

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                # Application entrypoint
├── internal/
│   ├── app/
│   │   ├── user/                  # User module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── dto.go
│   │   └── order/                 # Order module
│   │       ├── handler.go
│   │       ├── service.go
│   │       └── dto.go
│   ├── db/                        # sqlc generated code
│   ├── config/                    # Configuration
│   ├── middleware/                # HTTP middleware
│   ├── errors/                    # Custom errors
│   └── router/                    # Route registration
├── db/
│   ├── schema/                    # Database migrations
│   └── queries/                   # SQL query files
├── docker/
│   └── postgres/                  # Postgres init scripts
├── docker-compose.yml
├── sqlc.yaml
└── Makefile
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- golang-migrate CLI (for migrations)
- sqlc (for generating database code)

## Installation

### Install required tools

```bash
# Install golang-migrate
brew install golang-migrate

# Install sqlc
brew install sqlc

# Or using Go
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd bo
```

### 2. Create environment file

```bash
cp .env.example .env
```

### 2.5. Make sure Docker Desktop started

### 3. Start PostgreSQL with Docker

```bash
make docker-up
```

### 4. Run database migrations

```bash
make migrate-up
```

### 5. Generate sqlc code

```bash
make sqlc-generate
```

### 6. Install Go dependencies

```bash
go mod download
```

### 7. Run the application

```bash
make run
```

The API will be available at `http://localhost:8080`

## Makefile Commands

```bash
make help              # Show available commands
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make migrate-create    # Create new migration (usage: make migrate-create name=create_table)
make sqlc-generate     # Generate sqlc code
make swagger           # Generate Swagger documentation
make run               # Run the application
make build             # Build the application
make test              # Run tests
make clean             # Clean build artifacts
```

## API Documentation

Interactive API documentation is available via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

After modifying handlers or DTOs, regenerate the documentation:

```bash
make swagger
```

For detailed Swagger usage, see [docs/swagger-usage.md](docs/swagger-usage.md)

## API Endpoints

### Health Check

- `GET /health` - Health check endpoint

### Users

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - List users (supports pagination)
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Orders

- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - List orders (supports pagination and user filtering)
- `GET /api/v1/orders/{id}` - Get order by ID
- `PUT /api/v1/orders/{id}` - Update order
- `DELETE /api/v1/orders/{id}` - Delete order

## Example Requests

### Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Create Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "uuid-here",
    "order_number": "ORD-001",
    "status": "pending",
    "total_amount": 99.99,
    "notes": "First order"
  }'
```

### List Users with Pagination

```bash
curl "http://localhost:8080/api/v1/users?limit=10&offset=0"
```

### List Orders by User

```bash
curl "http://localhost:8080/api/v1/orders?user_id=uuid-here&limit=10&offset=0"
```

## Response Format

All API responses follow this format:

```json
{
  "status": true,
  "message": "descriptive message",
  "data": {
    // response data
  }
}
```

Error responses:

```json
{
  "status": false,
  "message": "error description"
}
```

## Development

### Adding a New Module

1. Create module directory: `internal/app/yourmodule/`
2. Create `dto.go`, `service.go`, and `handler.go`
3. Create SQL queries in `db/queries/yourmodule.sql`
4. Run `make sqlc-generate`
5. Register routes in `internal/router/router.go`

### Database Migrations

Create a new migration:

```bash
make migrate-create name=add_new_table
```

This creates two files in `db/schema/`:

- `{version}_add_new_table.up.sql`
- `{version}_add_new_table.down.sql`

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Production Deployment

1. Build the binary:

   ```bash
   make build
   ```

2. Set production environment variables:

   ```bash
   export DATABASE_URL="postgresql://..."
   export PORT="8080"
   ```

3. Run migrations:

   ```bash
   migrate -path db/schema -database $DATABASE_URL up
   ```

4. Start the application:
   ```bash
   ./bin/api
   ```

## Environment Variables

| Variable       | Description                  | Default                                                             |
| -------------- | ---------------------------- | ------------------------------------------------------------------- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable` |
| `PORT`         | HTTP server port             | `8080`                                                              |

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Chi v5
- **Database**: PostgreSQL 16
- **Database Driver**: pgx/v5
- **Query Builder**: sqlc
- **Validation**: go-playground/validator/v10
- **Logging**: slog (standard library)
- **Environment**: godotenv

## License

MIT

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
