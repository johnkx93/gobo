# GitHub Copilot Instructions

## Project Context

This is a Go-based REST API using PostgreSQL with Chi router, SQLC for type-safe queries, and JWT authentication.

## Core Principles

### 1. Audit Logging (CRITICAL)

- **ALWAYS** use `auditService` for tracking data changes (CREATE, UPDATE, DELETE operations)
- **ALWAYS** use `errorService` for logging application errors
- When implementing any CRUD operation:
  - After successful CREATE: call `auditService.LogCreate(ctx, tableName, entityID, newData)`
  - After successful UPDATE: call `auditService.LogUpdate(ctx, tableName, entityID, oldData, newData)`
  - After successful DELETE: call `auditService.LogDelete(ctx, tableName, entityID, oldData)`
- When catching errors that should be tracked: call `errorService.LogError(ctx, errorType, message, details)`
- Never skip audit logging - it's required for compliance and tracking

### 2. Database Migrations and Seeder Synchronization

- **ALWAYS** update `cmd/seeder/main.go` when:
  - Creating new migration files in `db/schema/`
  - Adding new tables or columns
  - Changing table structures
  - Adding new relationships or constraints
- The seeder must reflect the latest schema to generate valid test data
- Update the seeder's fake data generation to match new fields/tables
- Test the seeder after migration changes: `make seed`

### 3. SQLC Workflow

- After modifying SQL queries in `db/queries/`:
  - Run `sqlc generate` or `make sqlc-generate` immediately
  - Check for generated code in `internal/db/`
- Write SQL queries in `db/queries/*.sql` following SQLC annotations
- Never write raw SQL in Go code - use SQLC generated methods

### 4. Makefile Awareness

- **ALWAYS** check `Makefile` for existing commands before suggesting terminal commands
- When suggesting workflows, use Makefile targets:
  - `make dev` - Start full development environment
  - `make run` - Run the application
  - `make migrate-up` - Run migrations
  - `make sqlc-generate` - Generate SQLC code
  - `make seed` - Seed database with fake data
  - `make test` - Run tests
  - `make docker-up` / `make docker-down` - Manage Docker containers
- If a command doesn't exist in Makefile, suggest adding it there
- Reference Makefile targets in documentation and code comments

### 5. Error Handling Pattern

```go
// Correct pattern with audit and error logging
result, err := queries.SomeOperation(ctx, params)
if err != nil {
    // Log to error_logs table
    errorService.LogError(ctx, "DATABASE_ERROR", "Failed to perform operation", err.Error())
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Log successful audit trail
if err := auditService.LogCreate(ctx, "table_name", result.ID, result); err != nil {
    slog.Warn("failed to log audit", "error", err)
    // Continue - don't fail the operation if audit logging fails
}
```

### 6. Transaction Management for Complex Operations (CRITICAL)

- **ALWAYS** use transactions (`tx.Begin()` / `tx.Rollback()` / `tx.Commit()`) for complex operations involving multiple database writes
- This prevents **partial data insertion** when one operation fails mid-process
- Pattern for service methods with multiple database operations:

```go
func (s *Service) ComplexOperation(ctx context.Context, req Request) (*Response, error) {
    // Begin transaction
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        s.errorService.LogError(ctx, "TRANSACTION_ERROR", "Failed to begin transaction", err.Error())
        return nil, apperrors.NewInternalError("failed to start transaction")
    }
    defer tx.Rollback() // Always rollback on function exit (no-op if already committed)

    // Use transaction for all queries
    qtx := s.queries.WithTx(tx)

    // Step 1: First database operation
    result1, err := qtx.CreateEntity(ctx, params1)
    if err != nil {
        // Rollback happens automatically via defer
        s.errorService.LogError(ctx, "DATABASE_ERROR", "Failed to create entity", err.Error())
        return nil, apperrors.NewInternalError("failed to create entity")
    }

    // Step 2: Second database operation
    result2, err := qtx.UpdateRelatedEntity(ctx, params2)
    if err != nil {
        // Rollback happens automatically via defer
        s.errorService.LogError(ctx, "DATABASE_ERROR", "Failed to update related entity", err.Error())
        return nil, apperrors.NewInternalError("failed to update related entity")
    }

    // Step 3: Audit logging (also in transaction)
    if err := s.auditService.LogCreate(ctx, "entities", result1.ID, result1); err != nil {
        slog.Warn("failed to log audit", "error", err)
        // Continue - audit logging failure shouldn't fail the operation
    }

    // All operations successful - commit transaction
    if err := tx.Commit(); err != nil {
        s.errorService.LogError(ctx, "TRANSACTION_ERROR", "Failed to commit transaction", err.Error())
        return nil, apperrors.NewInternalError("failed to commit transaction")
    }

    return toResponse(result1, result2), nil
}
```

**When to use transactions:**

- ✅ Creating entity with related records (e.g., User + Address)
- ✅ Updating multiple tables that must stay consistent
- ✅ Deleting entity with cascading deletes
- ✅ Any operation where partial success would leave data inconsistent
- ❌ Single table INSERT/UPDATE/DELETE (not needed)
- ❌ Read-only queries (not needed)

**Key points:**

- `defer tx.Rollback()` ensures cleanup even if function panics
- `tx.Rollback()` after `tx.Commit()` is a no-op (safe to call)
- Use `queries.WithTx(tx)` to execute queries within the transaction
- **All or nothing** - either all operations succeed, or none do

### 7. Context Propagation

- Always pass `context.Context` through the call chain
- Audit context (user_id, request_id, ip_address, user_agent) is automatically available via middleware
- Use `audit.ExtractAuditContext(ctx)` when needed for manual operations

### 7. Service Layer Structure (CRITICAL)

#### Separate Services for Admin and User Operations

For modules that serve both admin and user functionality, **ALWAYS** create separate service classes:

**Structure:**

```
internal/app/entity/
├── dto.go              # Shared DTOs
├── admin_handler.go    # Admin handlers
├── admin_service.go    # AdminService - manages ANY entity
├── frontend_handler.go # User handlers
└── user_service.go     # UserService - enforces ownership
```

**Why Separate Services:**

- ✅ **No naming confusion** - both services can have same method names (e.g., `CreateEntity`)
- ✅ **Type safety** - compile-time enforcement prevents calling admin methods from user context
- ✅ **Clear intent** - handlers inject only the service they need
- ✅ **Better testability** - test admin and user logic independently
- ✅ **Easier to evolve** - admin and user features can diverge without conflicts

**AdminService Pattern:**

```go
type AdminService struct {
    queries      *db.Queries
    auditService *audit.Service
}

func NewAdminService(queries *db.Queries, auditService *audit.Service) *AdminService {
    return &AdminService{
        queries:      queries,
        auditService: auditService,
    }
}

// Admin can operate on ANY entity
func (s *AdminService) CreateEntity(ctx context.Context, req CreateEntityRequest) (*EntityResponse, error) {
    // Contains userID in request, admin specifies which user
    // Full access to all entities
}
```

**UserService Pattern:**

```go
type UserService struct {
    queries      *db.Queries
    auditService *audit.Service
}

func NewUserService(queries *db.Queries, auditService *audit.Service) *UserService {
    return &UserService{
        queries:      queries,
        auditService: auditService,
    }
}

// User can only operate on THEIR OWN entities
func (s *UserService) CreateEntity(ctx context.Context, userID uuid.UUID, req UserCreateEntityRequest) (*EntityResponse, error) {
    // userID is passed as parameter, enforced by service
    // Ownership validation built into queries
}
```

**Handler Injection:**

```go
// In cmd/api/main.go
entityAdminService := entity.NewAdminService(queries, auditService)
entityUserService := entity.NewUserService(queries, auditService)
entityAdminHandler := entity.NewAdminHandler(entityAdminService, validator)
entityFrontendHandler := entity.NewFrontendHandler(entityUserService, validator)
```

**When to Use Separate Services:**

- Module is used by both admin and regular users
- Different permission/ownership rules apply
- Business logic differs between admin and user operations

**When to Use Single Service:**

- Module is admin-only (e.g., `admin_management`)
- Module is user-only (e.g., `frontend_auth`)
- No ownership or permission differences

### 8. Migration Safety

- Always create both `.up.sql` and `.down.sql` files
- Test rollback before committing migrations
- For production migrations, remind about backup: `make migrate-up-prod` includes automatic backup
- Never force migrations in production without understanding the issue

### 9. Timestamp Columns (CRITICAL)

- **ALL tables** (except audit_logs and error_logs) MUST have:
  - `created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL`
  - `updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL`
- **Audit and error tables** only need:
  - `created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL`
- **Auto-update pattern**: Use PostgreSQL triggers to automatically update `updated_at`

#### Standard Trigger Pattern for `updated_at`:

```sql
-- Create reusable trigger function (do this once in a migration)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to each table (in the table's migration)
CREATE TRIGGER trigger_update_TABLENAME_updated_at
    BEFORE UPDATE ON TABLENAME
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**Note:** The `update_updated_at_column()` function already exists in the database (created in migration 008). For new tables, you only need to create the trigger, not the function.

#### Example for New Tables:

```sql
-- Create table with timestamps
CREATE TABLE new_entity (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Add trigger (the function already exists, just reference it)
CREATE TRIGGER trigger_update_new_entity_updated_at
    BEFORE UPDATE ON new_entity
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### 10. Testing

- Write tests for new services and handlers
- Run tests before committing: `make test`
- For coverage reports: `make test-coverage`

#### Unit Testing Guidelines (On-Demand)

Unit tests are created when requested, not automatically. Follow this two-step process:

**Step 1: Service Layer Unit Tests (Validation Focus)**

1. **Analyze Service Layer:**

   - Identify all public methods in service files
   - Understand validation rules and business logic
   - Note dependencies (database, audit service)

2. **Identify Unit-Testable Logic (No DB Required):**

   - ✅ Input validation (UUID parsing, required fields)
   - ✅ Helper functions (getStringValue, toResponse converters)
   - ✅ Error handling (correct error types/codes)
   - ✅ Business logic without DB
   - ❌ Database operations (save for integration tests)
   - ❌ Audit logging (save for integration tests)

3. **Create Service Test Files:**

   - Follow naming: `admin_service.go` → `admin_service_test.go`
   - Same package as the code being tested
   - Structure:

     ```go
     package modulename

     import (
         "context"
         "testing"
         // imports
     )

     // 1. Test helper functions
     // 2. Test service methods (validation focus)
     ```

4. **Write Validation-Focused Tests:**

   - Test validation logic by setting dependencies to `nil`:
     ```go
     service := &AdminService{
         queries: nil,      // Will panic if DB called
         auditService: nil, // Only testing validation
     }
     ```
   - Use table-driven tests for multiple cases
   - Pattern:

     ```go
     func TestServiceMethod_Scenario(t *testing.T) {
         // Arrange
         service := &Service{queries: nil, auditService: nil}
         ctx := context.Background()

         // Act
         _, err := service.Method(ctx, invalidInput)

         // Assert
         if err == nil {
             t.Fatal("expected error, got nil")
         }

         domainErr, ok := err.(*apperrors.DomainError)
         if !ok {
             t.Fatalf("expected DomainError, got %T", err)
         }

         if domainErr.Code != apperrors.CodeValidation {
             t.Errorf("expected validation error, got: %s", domainErr.Code)
         }
     }
     ```

5. **Run Service Tests:**
   ```bash
   go test -v ./internal/app/modulename/... -run "TestAdminService|TestFrontendService"
   ```

**Step 2: Handler Layer Unit Tests (HTTP Focus)**

1. **Analyze Handler Layer:**

   - Identify all handler methods (admin and frontend)
   - Understand authentication requirements (admin role, user ID)
   - Note request/response patterns

2. **Create Handler Test Files:**

   - Follow naming: `admin_handler.go` → `admin_handler_test.go`
   - Same package as the code being tested
   - Structure:

     ```go
     package modulename

     import (
         "bytes"
         "context"
         "encoding/json"
         "net/http"
         "net/http/httptest"
         "testing"

         "github.com/go-chi/chi/v5"
         "github.com/user/coc/internal/ctxkeys"
         "github.com/user/coc/internal/validation"
     )

     // 1. Mock service interfaces
     // 2. Helper functions (newAdminRequest, newUserRequest)
     // 3. Test handler methods
     ```

3. **Create Mock Services:**

   - Define mock interfaces that return controlled responses:

     ```go
     type MockAdminService struct {
         CreateFunc func(context.Context, CreateEntityRequest) (*EntityResponse, error)
         GetFunc    func(context.Context, string) (*EntityResponse, error)
         // ... other methods
     }

     func (m *MockAdminService) CreateEntity(ctx context.Context, req CreateEntityRequest) (*EntityResponse, error) {
         if m.CreateFunc != nil {
             return m.CreateFunc(ctx, req)
         }
         return nil, nil
     }
     ```

4. **Create Request Helper Functions:**

   - For admin handlers:

     ```go
     func newAdminRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
         var reqBody []byte
         if body != nil {
             reqBody, _ = json.Marshal(body)
         }
         req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
         req.Header.Set("Content-Type", "application/json")

         // Add admin role to context
         ctx := context.WithValue(req.Context(), ctxkeys.AdminRoleContextKey, "super_admin")
         req = req.WithContext(ctx)

         rec := httptest.NewRecorder()
         return req, rec
     }
     ```

   - For frontend handlers:

     ```go
     func newUserRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
         var reqBody []byte
         if body != nil {
             reqBody, _ = json.Marshal(body)
         }
         req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
         req.Header.Set("Content-Type", "application/json")

         // Add user ID to context
         userID := uuid.New()
         ctx := context.WithValue(req.Context(), ctxkeys.UserIDContextKey, userID)
         req = req.WithContext(ctx)

         rec := httptest.NewRecorder()
         return req, rec
     }
     ```

5. **Write Handler Tests (Focus Areas):**

   - **Authentication Tests:**

     ```go
     func TestAdminHandler_CreateEntity_MissingAdminRole(t *testing.T) {
         handler := NewAdminHandler(&MockAdminService{}, validation.New())
         req := httptest.NewRequest("POST", "/entities", nil)
         // No admin role in context
         rec := httptest.NewRecorder()

         handler.CreateEntity(rec, req)

         if rec.Code != http.StatusUnauthorized {
             t.Errorf("expected 401, got %d", rec.Code)
         }
     }
     ```

   - **Missing Parameters Tests:**

     ```go
     func TestFrontendHandler_GetEntity_MissingEntityID(t *testing.T) {
         handler := NewFrontendHandler(&MockFrontendService{}, validation.New())
         req, rec := newUserRequest("GET", "/entities/", nil)
         // Missing entity ID in URL params

         handler.GetEntity(rec, req)

         if rec.Code != http.StatusBadRequest {
             t.Errorf("expected 400, got %d", rec.Code)
         }
     }
     ```

   - **Invalid JSON Tests:**

     ```go
     func TestAdminHandler_CreateEntity_InvalidJSON(t *testing.T) {
         handler := NewAdminHandler(&MockAdminService{}, validation.New())
         req, rec := newAdminRequest("POST", "/entities", nil)
         req.Body = io.NopCloser(bytes.NewBufferString("invalid json"))

         handler.CreateEntity(rec, req)

         if rec.Code != http.StatusBadRequest {
             t.Errorf("expected 400, got %d", rec.Code)
         }
     }
     ```

   - **Validation Error Tests:**

     ```go
     func TestAdminHandler_CreateEntity_ValidationError(t *testing.T) {
         handler := NewAdminHandler(&MockAdminService{}, validation.New())
         req, rec := newAdminRequest("POST", "/entities", map[string]interface{}{
             "user_id": "invalid-uuid",
         })

         handler.CreateEntity(rec, req)

         if rec.Code != http.StatusBadRequest {
             t.Errorf("expected 400, got %d", rec.Code)
         }
     }
     ```

6. **Run Handler Tests:**
   ```bash
   go test -v ./internal/app/modulename/... -run "TestAdminHandler|TestFrontendHandler"
   ```

**Run All Unit Tests:**

```bash
go test -v ./internal/app/modulename/...
go test -v ./internal/app/modulename/... -run TestSpecific
go test -cover ./internal/app/modulename/...
```

**Test File Structure Example:**

```
internal/app/address/
├── admin_service.go
├── admin_service_test.go       # Step 1: Service validation tests
├── admin_handler.go
├── admin_handler_test.go       # Step 2: Handler HTTP tests
├── frontend_service.go
├── frontend_service_test.go    # Step 1: Service validation tests
├── frontend_handler.go
├── frontend_handler_test.go    # Step 2: Handler HTTP tests
├── integration_test.go         # Integration tests with real DB (future)
└── dto.go
```

**What Unit Tests Should Cover:**

**Service Tests:**

- Invalid UUID formats
- Missing required fields (if validated in service)
- Helper function edge cases (nil pointers, empty strings)
- Error type and code correctness

**Handler Tests:**

- Missing authentication context (admin role, user ID)
- Missing URL parameters (entity IDs)
- Invalid JSON payloads
- Validation errors (invalid UUIDs, missing fields)
- HTTP status codes (200, 201, 400, 401, 404)
- Response structure (status, message, data fields)

**What Unit Tests Should NOT Cover (Use Integration Tests):**

- Database CRUD operations
- Audit trail creation
- Complex queries and joins
- Transaction handling
- Trigger functionality
- Full request-to-database workflows

**Step 3: Integration Tests (Real Database)**

Integration tests verify the complete flow from service to database using real PostgreSQL.

1. **Create Integration Test File:**

   - File naming: `integration_test.go` in the module folder
   - Uses real Docker Postgres database
   - All tests use transactions with rollback for cleanup

2. **Setup Test Database Helper:**

   ```go
   func setupTestDB(t *testing.T) (*pgxpool.Pool, *db.Queries, *audit.Service) {
       t.Helper()

       // Use DATABASE_URL from environment (Docker Postgres)
       databaseURL := os.Getenv("DATABASE_URL")
       if databaseURL == "" {
           t.Skip("DATABASE_URL not set, skipping integration tests")
       }

       ctx := context.Background()

       // Parse and create pool config
       config, err := pgxpool.ParseConfig(databaseURL)
       if err != nil {
           t.Fatalf("failed to parse database URL: %v", err)
       }

       // Configure pool for testing
       config.MaxConns = 10
       config.MinConns = 2

       pool, err := pgxpool.NewWithConfig(ctx, config)
       if err != nil {
           t.Fatalf("failed to create connection pool: %v", err)
       }

       // Verify connection
       if err := pool.Ping(ctx); err != nil {
           pool.Close()
           t.Fatalf("failed to ping database: %v", err)
       }

       // Create queries and services
       queries := db.New(pool)
       auditService := audit.NewService(queries)

       // Cleanup on test completion
       t.Cleanup(func() {
           pool.Close()
       })

       return pool, queries, auditService
   }
   ```

3. **Write Integration Tests with Transactions:**

   ```go
   func TestIntegration_AdminService_CreateEntity(t *testing.T) {
       pool, queries, auditService := setupTestDB(t)

       ctx := context.Background()

       // Start transaction for isolation
       tx, err := pool.Begin(ctx)
       if err != nil {
           t.Fatalf("failed to begin transaction: %v", err)
       }
       defer tx.Rollback(ctx) // CRITICAL: Cleanup test data

       // Use transaction for queries
       qtx := queries.WithTx(tx)
       service := NewAdminService(qtx, auditService)

       // Create test data
       result, err := service.CreateEntity(ctx, createRequest)
       if err != nil {
           t.Fatalf("CreateEntity failed: %v", err)
       }

       // Verify data in database
       dbEntity, err := qtx.GetEntity(ctx, entityID)
       if err != nil {
           t.Fatalf("failed to get entity from database: %v", err)
       }

       // Verify audit log created
       auditLogs, err := qtx.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
           EntityType: "entities",
           EntityID:   pgtype.UUID{Bytes: entityID, Valid: true},
           Limit:      10,
           Offset:     0,
       })
       if err != nil {
           t.Fatalf("failed to get audit logs: %v", err)
       }

       if len(auditLogs) == 0 {
           t.Error("expected audit log entry, got none")
       }

       if auditLogs[0].Action != "CREATE" {
           t.Errorf("expected audit action CREATE, got %s", auditLogs[0].Action)
       }

       // Transaction rolls back automatically - no cleanup needed!
   }
   ```

4. **What Integration Tests Should Cover:**

   - ✅ Full CRUD workflows (Create → Read → Update → Delete)
   - ✅ Database operations with SQLC queries
   - ✅ Audit trail creation in `audit_logs` table
   - ✅ Complex queries and joins
   - ✅ Transaction handling
   - ✅ Trigger functionality (like `updated_at` auto-update)
   - ✅ Foreign key constraints and cascades
   - ✅ Data integrity and validation at DB level
   - ✅ Ownership enforcement (user-owned resources)
   - ✅ Error handling for duplicate entries, constraint violations

5. **Run Integration Tests:**

   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable"
   go test -v ./internal/app/modulename/... -run "TestIntegration"
   ```

6. **Key Integration Test Patterns:**
   - **Transaction Isolation:** Each test starts with `tx.Begin()` and ends with `defer tx.Rollback(ctx)`
   - **No Leftover Data:** Rollback ensures database stays clean
   - **Real Database:** Tests actual SQL queries, triggers, and constraints
   - **Skip if No DB:** Tests skip gracefully if `DATABASE_URL` not set
   - **Helper Functions:** Create test users, entities with `createTestUser()`, etc.

**Complete Test File Structure:**

```
internal/app/address/
├── admin_service.go
├── admin_service_test.go       # Step 1: Service validation tests (no DB)
├── admin_handler.go
├── admin_handler_test.go       # Step 2: Handler HTTP tests (no DB)
├── frontend_service.go
├── frontend_service_test.go    # Step 1: Service validation tests (no DB)
├── frontend_handler.go
├── frontend_handler_test.go    # Step 2: Handler HTTP tests (no DB)
├── integration_test.go         # Step 3: Integration tests (real DB)
└── dto.go
```

**Running All Tests:**

```bash
# Unit tests only (fast, no database required)
go test -v ./internal/app/modulename/... -short

# Integration tests only (requires DATABASE_URL)
export DATABASE_URL="postgres://..."
go test -v ./internal/app/modulename/... -run "TestIntegration"

# All tests (unit + integration)
export DATABASE_URL="postgres://..."
go test -v ./internal/app/modulename/...

# With coverage
go test -cover ./internal/app/modulename/...
```

### 11. Security Best Practices

- Audit service automatically filters sensitive fields (password_hash, token, secret, etc.)
- Never log passwords or tokens in plain text
- Always hash passwords using bcrypt before storing
- JWT tokens should have appropriate expiration times

### 12. Admin Handler Pattern (CRITICAL)

- **ALWAYS** check admin role at the start of admin handler functions
- Use `ctxkeys.GetAdminRole(r)` to retrieve the role from context
- Return `401 Unauthorized` if role is not found or empty
- Pattern to follow:

```go
func (h *Handler) SomeAdminAction(w http.ResponseWriter, r *http.Request) {
    // REQUIRED: Check admin role first
    role, ok := ctxkeys.GetAdminRole(r)
    if !ok || role == "" {
        response.Error(w, http.StatusUnauthorized, "admin role not found")
        return
    }

    // Continue with handler logic...
}
```

- This prevents unauthorized access even if middleware is misconfigured
- All handlers in `internal/app/admin_*` packages must follow this pattern
- Use `ctxkeys.GetAdminID(r)` when you need the admin's UUID for audit logging

### 13. Permissions and Menu System (CRITICAL)

- **ALWAYS** add permissions and menu items when creating new tables/modules
- Three tables exist for access control:
  - `permissions` - Defines permission codes, names, and categories
  - `role_permissions` - Maps roles to permissions (many-to-many)
  - `menu_items` - Defines admin menu structure with permission requirements

#### Adding New Module Permissions:

```sql
-- 1. Add permissions for the new module
INSERT INTO permissions (code, name, description, category) VALUES
    ('module.create', 'Create Module', 'Description', 'module'),
    ('module.read', 'Read Module', 'Description', 'module'),
    ('module.update', 'Update Module', 'Description', 'module'),
    ('module.delete', 'Delete Module', 'Description', 'module');

-- 2. Assign to super_admin role (gets ALL permissions)
INSERT INTO role_permissions (role, permission_id)
SELECT 'super_admin', id FROM permissions WHERE category = 'module';

-- 3. Assign to admin role (usually read/update only)
INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions WHERE code IN ('module.read', 'module.update');

-- 4. Assign to moderator role (usually read-only)
INSERT INTO role_permissions (role, permission_id)
SELECT 'moderator', id FROM permissions WHERE code = 'module.read';
```

#### Adding Menu Items:

```sql
-- 1. Add root menu item
INSERT INTO menu_items (code, label, icon, order_index, permission_id) VALUES
    ('module', 'Module Management', 'icon-name', 6, NULL);

-- 2. Add child menu items
INSERT INTO menu_items (parent_id, code, label, path, order_index, permission_id)
SELECT
    (SELECT id FROM menu_items WHERE code = 'module'),
    'module-list',
    'Module List',
    '/admin/modules',
    1,
    (SELECT id FROM permissions WHERE code = 'module.read')
UNION ALL
SELECT
    (SELECT id FROM menu_items WHERE code = 'module'),
    'module-create',
    'Create Module',
    '/admin/modules/create',
    2,
    (SELECT id FROM permissions WHERE code = 'module.create');
```

### 14. Database Access (CRITICAL)

- **ALWAYS** use Docker to access PostgreSQL - never direct psql
- Pattern for database queries:

```bash
docker exec -it $(docker ps -q -f name=postgres) psql -U postgres -d appdb -c "YOUR_QUERY"
```

- This ensures consistency across development environments
- Works on all platforms (macOS, Linux, Windows with WSL)
- Prevents connection issues and environment variable conflicts

## Common Patterns

### Adding a New Entity

1. Create migration files: `make migrate-create name=create_entity_table`
2. Write SQL schema in `db/schema/XXX_create_entity_table.up.sql`:
   - Include `created_at` and `updated_at` columns
   - Add trigger for auto-updating `updated_at`
3. Add SQLC queries in `db/queries/entity.sql`
4. Run migrations: `make migrate-up`
5. Generate SQLC code: `make sqlc-generate`
6. **Create application structure in `internal/app/entity/`**:

   **For Shared Modules (both admin and user):**

   - `dto.go` - Request/response DTOs with validation tags
     - Admin DTOs: `CreateEntityRequest`, `UpdateEntityRequest` (contain userID)
     - User DTOs: `UserCreateEntityRequest`, `UserUpdateEntityRequest` (no userID field)
     - Shared: `EntityResponse`
   - `admin_service.go` - Admin business logic (manages ANY entity)
   - `user_service.go` - User business logic (enforces ownership)
   - `admin_handler.go` - Admin HTTP handlers
   - `frontend_handler.go` - User HTTP handlers

   **For Single-Party Modules:**

   - `dto.go` - Request/response DTOs
   - `service.go` - Business logic with audit logging
   - `handler.go` - HTTP handlers

7. **Add routes in `internal/router/`**:
   - Admin routes → `admin.go` (protected by admin auth middleware)
   - User routes → `frontend.go` (protected by user auth middleware)
   - Register handlers and apply appropriate middleware
8. **Update seeder in `cmd/seeder/main.go`** to generate fake data for new entity
9. **CRITICAL: Add permissions and menu items** (create new migration):
   - Add permissions to `permissions` table (category = entity name)
   - Assign permissions to roles in `role_permissions` (super_admin gets all)
   - Add menu items to `menu_items` table with proper permission links
   - Update `sqlc.yaml` to include the new migration file

### Updating Existing Entity

1. Create new migration: `make migrate-create name=add_field_to_entity`
2. Write ALTER TABLE statements
3. Run migration: `make migrate-up`
4. Update SQLC queries if needed and regenerate
5. Update service methods to include new fields
6. Update DTOs for validation
7. **Update seeder to populate new fields**
8. Update existing audit logs to capture new fields

## File Organization

### Naming Conventions (CRITICAL)

#### File Names:

- **Admin handlers:** `admin_handler.go` (NOT `handler_admin.go`)
- **frontend handlers:** `frontend_handler.go` (NOT `handler_frontend.go`)
- **Admin services:** `admin_service.go` (NOT `service_admin.go`)
- **Frontend services:** `frontend_service.go` (NOT `service_frontend.go`)
- **DTOs:** `dto.go`
- **Single handler:** `handler.go` (for modules with only one user type)
- **Single service:** `service.go` (for modules with only one user type)

**Examples:**

```
✅ CORRECT:
internal/app/order/
├── admin_handler.go
├── admin_service.go
├── frontend_handler.go
├── frontend_service.go
└── dto.go

❌ INCORRECT:
internal/app/order/
├── handler_admin.go
├── service_admin.go
├── handler_frontend.go
├── service_user.go
└── dto.go
```

#### Struct Names:

- **Admin handlers:** `AdminHandler` (injected with `AdminService`)
- **Frontend handlers:** `FrontendHandler` (injected with `FrontendService`)
- **Admin services:** `AdminService`
- **Frontend services:** `FrontendService`

#### Method Names:

**Admin Service Methods** (no prefix needed):

```go
// AdminService methods - operate on ANY entity
func (s *AdminService) CreateEntity(...)
func (s *AdminService) GetEntity(...)
func (s *AdminService) ListEntities(...)
func (s *AdminService) UpdateEntity(...)
func (s *AdminService) DeleteEntity(...)
```

**Frontend Service Methods** (no prefix needed):

```go
// FrontendService methods - enforce ownership
func (s *FrontendService) CreateEntity(ctx, userID uuid.UUID, ...)
func (s *FrontendService) GetEntity(ctx, userID uuid.UUID, entityID, ...)
func (s *FrontendService) ListEntities(ctx, userID uuid.UUID, ...)
func (s *FrontendService) UpdateEntity(ctx, userID uuid.UUID, entityID, ...)
func (s *FrontendService) DeleteEntity(ctx, userID uuid.UUID, entityID, ...)
```

**Handler Methods** (match service method names):

```go
// AdminHandler
func (h *AdminHandler) CreateEntity(w, r) { ... }

// FrontendHandler
func (h *FrontendHandler) CreateEntity(w, r) { ... }
```

#### DTO Names:

**Admin DTOs** (contain userID or target user):

```go
type CreateEntityRequest struct {
    UserID string `json:"user_id" validate:"required,uuid"`
    Name   string `json:"name" validate:"required"`
}

type UpdateEntityRequest struct {
    Name *string `json:"name" validate:"omitempty"`
}
```

**Frontend DTOs** (prefix with "Frontend", no userID field):

```go
type FrontendCreateEntityRequest struct {
    // NO UserID field - enforced by service
    Name string `json:"name" validate:"required"`
}

type FrontendUpdateEntityRequest struct {
    Name *string `json:"name" validate:"omitempty"`
}
```

**Shared DTOs:**

```go
type EntityResponse struct {
    ID        string `json:"id"`
    UserID    string `json:"user_id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

### Directory Structure

- `cmd/` - Entry points (main.go, seeder)
- `internal/app/` - Business logic by domain (user, order, auth)

  - Each module folder structure depends on usage:

  **For Shared Modules (both admin and frontend):**

  ```
  internal/app/entity/
  ├── dto.go              # Shared DTOs (both admin and frontend)
  ├── admin_handler.go    # Admin HTTP handlers
  ├── admin_service.go    # Admin business logic
  ├── frontend_handler.go # Frontend HTTP handlers
  └── frontend_service.go # Frontend business logic with ownership checks
  ```

  **For Admin-Only Modules:**

  ```
  internal/app/admin_management/
  ├── dto.go
  ├── handler.go
  └── service.go
  ```

  **For Frontend-Only Modules:**

  ```
  internal/app/frontend_auth/
  ├── dto.go
  ├── handler.go
  └── service.go
  ```

- `internal/audit/` - Audit and error logging services
- `internal/db/` - SQLC generated code (don't edit manually)
- `internal/middleware/` - HTTP middleware
- `internal/router/` - Route definitions
  - `router.go` - Main router setup
  - `admin_router.go` - Admin routes (protected by admin auth middleware)
  - `frontend_router.go` - Frontend routes (protected by frontend auth middleware)
- `db/schema/` - Migration files
- `db/queries/` - SQLC query definitions
- `scripts/` - Shell scripts for operations
- `docs/` - Documentation

## Remember

- Audit logging is not optional - it's a core feature
- Always sync seeder with migrations
- Use Makefile commands consistently
- Test thoroughly before suggesting production changes
- **ALWAYS use Docker for PostgreSQL access** - never direct `psql`
- **Every new table/module MUST have permissions and menu items** - this is mandatory for admin access control
