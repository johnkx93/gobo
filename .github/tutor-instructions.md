## Tutor instructions — Backend Order (BO) project

Purpose
-------

This document is written for an AI tutor that will teach a student how this project works. The goal is to explain the overall architecture and the runtime flow, point the student to the most important files, provide step-by-step commands to run and test the app locally, and propose small exercises and debugging tips so the student can learn by doing.

Use this file as the canonical guide when you (the tutor) are walking a student through the codebase.

Quick orientation
-----------------

- Language: Go (>= 1.21)
- DB: PostgreSQL (Docker-supported)
- Router: Chi
- DB codegen: sqlc
- Migrations: golang-migrate

Top-level commands (Makefile)
- Start Docker DB: `make docker-up`
- Run migrations: `make migrate-up`
- Generate sqlc code: `make sqlc-generate`
- Run the app: `make run`
- Start dev environment (docker + migrate + sqlc + run): `make dev`
- Setup (deps + docker + migrate + sqlc): `make setup`

Quickstart (what to do with a fresh clone)
-----------------------------------------
1. Ensure prerequisites: Go, Docker, Docker Compose, `sqlc`, and `migrate` CLI tools.
2. Start the DB: `make docker-up`
3. Run migrations: `make migrate-up`
4. Generate typed DB code: `make sqlc-generate`
5. Run the app: `make run`

The API should be reachable at `http://localhost:8080` by default.

Where to look first (a short guided walkthrough)
-----------------------------------------------

1. `cmd/api/main.go` — application entrypoint and wiring: logger, config, DB pool, services, handlers, router, and graceful shutdown. This is the best single-file overview of how pieces are connected.

2. `internal/config/config.go` — how configuration is loaded (env vars, defaults). Look for `Port` and `DatabaseURL` variables used on startup.

3. `internal/db/pool.go` and `internal/db/querier.go` (or `db.go`) — database pool creation and the `New(pool)` function that returns the typed queries from sqlc.

4. `internal/router/router.go` — route registration and middleware. It shows the HTTP paths and which handlers are mounted.

5. Modules: `internal/app/user` and `internal/app/order` — each module typically contains three files:
	- `dto.go` — request/response shapes and validation tags
	- `handler.go` — HTTP layer: decode request, validate, call service, format response
	- `service.go` — business logic: transactional rules, calling sqlc-generated queries

6. `db/queries/*.sql` and `internal/db/*.go` — sqlc inputs and generated outputs. SQL files live under `db/queries/` and the generated typed code lives in `internal/db/`.

Data flow for a typical request (Create Order request example)
-----------------------------------------------------------

1. Client sends POST /api/v1/orders with JSON body.
2. Chi router matches the route and invokes the handler from `internal/app/order/handler.go`.
3. Handler decodes JSON into a DTO defined in `dto.go` and runs validation (validator tags).
4. Handler calls the `order.Service` method (in `service.go`) with the DTO or converted domain object.
5. Service contains business rules and uses the sqlc-generated `queries` methods (from `internal/db`) to insert/select/update records.
6. sqlc code executes SQL against the connection pool created in `internal/db/pool.go`.
7. Service returns the domain result to the handler, which marshals JSON and responds.

Key files and their purpose (cheat sheet)
----------------------------------------

- `cmd/api/main.go` — app bootstrap and graceful shutdown
- `internal/config/config.go` — config loading and defaults
- `internal/db/pool.go` — pgx-based connection pool
- `internal/db/querier.go`, `internal/db/*.sql.go` — sqlc-generated queries and helper wrapper
- `internal/router/router.go` — HTTP route registration and middleware
- `internal/middleware/*` — middleware (logging, common hooks)
- `internal/app/*/handler.go` — HTTP handlers for each module
- `internal/app/*/service.go` — business logic
- `db/schema/` — migrations (`*.up.sql` and `*.down.sql`)
- `db/queries/` — SQL files used by sqlc
- `docker/` — DB init scripts (e.g., `docker/postgres/init.sql`)

Database migrations and docker
------------------------------

- Migrations live in `db/schema/`. New migrations are created with `make migrate-create name=your_name`.
- Run migrations with `make migrate-up` (Makefile uses the `migrate` CLI and targets the default connection string). For production, prefer passing a `DATABASE_URL`.
- Docker Compose contains a Postgres service used for local development. Start it with `make docker-up` and stop with `make docker-down`.

Configuration and environment
------------------------------

- Default DB connection string used in Makefile: `postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable`.
- Typical env vars to set/test:
  - `DATABASE_URL` — connection string
  - `PORT` — server port (defaults to 8080)

Testing and code generation
---------------------------

- Generate DB client code: `make sqlc-generate` (runs `sqlc generate` using `sqlc.yaml`).
- Run unit/integration tests: `make test` or `go test ./...`.

Suggested learning exercises (progressive)
-----------------------------------------

1. Trace a single request start-to-finish
	- Run the app locally.
	- Use `curl` to POST a user, then GET the user.
	- Open `main.go`, `router.go`, `handler.go`, `service.go`, and the generated `internal/db/*.go` for the query that ran.

2. Add a small endpoint
	- Add a `GET /api/v1/users/me` that returns a fixed mock response.
	- Add route in `router.go`, small handler in `internal/app/user/handler.go`, and a test for it.

3. Add validation
	- Modify a DTO to add or adjust validation tags and observe request rejections.

4. Add a migration and update sqlc
	- Create a migration with `make migrate-create name=add_thing`.
	- Add SQL in `db/queries/` for new behavior and run `make sqlc-generate`.

5. Write a unit test for a service method
	- Mock `queries` or use a test DB and assert expected behavior.

Debugging tips
--------------

- Check Docker container logs: `make docker-logs` or `docker-compose logs -f`.
- Inspect Postgres directly: `docker exec -it <postgres_container> psql -U postgres -d appdb`.
- The app uses structured JSON logging (slog). Read logs to see request lifecycle and errors.
- If SQL queries from sqlc fail, open the generated file in `internal/db/` to see the exact SQL and parameters.
- For panic or startup issues, inspect `cmd/api/main.go` order: config -> db.NewPool -> db.New -> service constructors. A failure in any of these aborts startup.

Tutor prompts (use these when walking the student through the code)
--------------------------------------------------------------

1. "Show me how the server starts — what happens before the HTTP router is created?"
2. "Where is the DB connection created and how is it passed to the code that runs queries?"
3. "Find the code that creates an order — show me handler -> service -> sqlc call." 
4. "Add an assertion: what should happen if validation fails in a DTO?"

Common pitfalls and edge cases for students to watch
--------------------------------------------------

- Not regenerating sqlc code after changing SQL files — always run `make sqlc-generate`.
- Forgetting to run migrations when DB schema changes.
- Mixing environment DB URLs — Makefile uses localhost; Docker can map ports; ensure `DATABASE_URL` matches.
- Tests that depend on a real DB: prefer using a test container or mocking queries.

Suggested next steps for the tutor
---------------------------------

1. Ask the student to run the quickstart and create a user and an order. Watch logs together.
2. Have the student open `cmd/api/main.go` and draw (or describe) the wiring: config -> db -> queries -> services -> handlers -> router.
3. Assign the "Add a new endpoint" exercise and review the PR together.

References
----------

- Migrations: `db/schema/`
- SQL queries: `db/queries/`
- Generated DB client: `internal/db/`
- App entrypoint: `cmd/api/main.go`
- Modules: `internal/app/user` and `internal/app/order`

Completion note
---------------

This file is intentionally hands-on: include concrete exercises and short walkthrough prompts for the student. If you want, I can also generate a short checklist of 5 paired-programming exercises with step-by-step hints and starter patches.

