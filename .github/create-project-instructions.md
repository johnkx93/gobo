# GitHub Copilot Instructions

You are GitHub Copilot assisting in a **Go backend** that uses:

- **PostgreSQL** (database via Docker)
- **Go + pgx + sqlc** (type-safe queries)
- **Modular architecture**: each feature (user, order, etc.) is an independent module with its own handlers, service logic, and SQL queries.

Your goals are:
- Keep Docker setup simple and reliable for Postgres.
- Maintain a clean modular folder structure.
- Generate Go code that uses pgx + sqlc efficiently.
- Keep business logic isolated per module.

---

## 1. Project structure expectations

Overall layout:

```text
.
├── cmd/
│   └── api/
│       └── main.go                # app entrypoint
│
├── internal/
│   ├── app/
│   │   ├── user/                  # user module
│   │   │   ├── handler.go         # HTTP layer
│   │   │   ├── service.go         # business logic
│   │   │   └── dto.go             # request/response structs
│   │   ├── order/                 # order module
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── dto.go
│   │   └── ...                    # other modules
│   │
│   ├── db/                        # sqlc-generated code (package "db")
│   ├── config/                    # env & app configuration
│   ├── middleware/                # shared middlewares (auth, logging, etc.)
│   └── router/                    # route registration per module
│
├── db/
│   ├── schema/                    # schema & migrations
│   └── queries/                   # sqlc query files
│       ├── user.sql
│       ├── order.sql
│       └── ...
│
├── docker/
│   └── postgres/                  # optional init scripts / config
│
├── docker-compose.yml
├── sqlc.yaml
└── copilot-instructions.md
```

---

## 2. Database conventions
- Use snake_case for tables and columns
- Migrations in `db/schema/` numbered sequentially
- Use sqlc for all database queries
- Database driver: `pgx/v5`
- Connection pooling: `pgxpool`
- Type-safe queries: `sqlc` generated into `internal/db`
- No ORM (avoid GORM or raw SQL unless necessary)

---

## 3. HTTP layer patterns
- Framework: [chi]
- Handler signature: func(w http.ResponseWriter, r *http.Request)
- Use DTOs for request/response validation
- Return JSON responses with consistent structure
- Validation library: use github.com/go-playground/validator/v10 with struct tags
- Handlers are responsible for:
	- Parsing input (path/query/body).
	- Validating DTOs (using validator).
	- Calling service methods.
	- Mapping service errors to HTTP status codes and the JSON error envelope.
- JSON response envelope:
```json
{
    "status": true | false,
    "message": "descriptive message",
    "data": {...},
}
```

---

## 4. Error handling
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use custom error types in `internal/errors/` for domain errors
- Log errors at service layer using [slog]
- Return HTTP status codes: 400 (bad request), 404 (not found), 500 (internal error)
- Never expose internal error details to clients

---

## 5. Environment & configuration
- Use .env for local development
- Required vars: DATABASE_URL, PORT and any feature-specific variables.

---

## 6. Use Docker for local development:
- Service name: `postgres`
- Image: `postgres:16-alpine`
- Port: 5432
- Default database: `appdb`

---

## 7. Coding style
- Use `context.Context` in all database and service methods
- Return `(T, error)` instead of panicking
- Keep functions small and focused
- Only import minimal dependencies (prefer stdlib + pgx + sqlc)

---

## 8. Testing
- Unit tests alongside code: `handler_test.go`, `service_test.go`
- Use testcontainers for integration tests
- Mock database with interfaces
- Test coverage target: 80%+

---

## 9. Logging
- Use structured logging with slog
- Log levels: debug, info, warn, error
- Include request IDs in logs

---

## 10. Database migrations
- Tool: golang-migrate/migrate
- Format: `{version}_{name}.up.sql` and `{version}_{name}.down.sql`
- Run via: `migrate -path db/schema -database $DATABASE_URL up`

---

## 11. sqlc.yaml configuration
- Package: `db`
- Output: `internal/db`
- Engine: `postgresql`
- Emit: json_tags, prepared_queries, interface