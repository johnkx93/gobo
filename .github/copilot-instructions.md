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

### 6. Context Propagation
- Always pass `context.Context` through the call chain
- Audit context (user_id, request_id, ip_address, user_agent) is automatically available via middleware
- Use `audit.ExtractAuditContext(ctx)` when needed for manual operations

### 7. Service Layer Structure
- Services should be injected with:
  - `queries *db.Queries` - Database queries
  - `auditService *audit.Service` - Audit logging
  - `errorService *audit.ErrorService` - Error logging (if needed)
- Follow existing patterns in `internal/app/user/service.go` and `internal/app/order/service.go`

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
   - `dto.go` - Request/response DTOs and validation
   - `service.go` - Business logic with audit logging
   - **If module is used by both admin and users:**
     - `handler_admin.go` - Admin-specific handlers
     - `handler_frontend.go` - User-facing handlers (frontend)
   - **If module is used by only one party:**
     - `handler.go` - Single handler file
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
- `cmd/` - Entry points (main.go, seeder)
- `internal/app/` - Business logic by domain (user, order, auth)
  - Each module folder contains:
    - `dto.go` - Request/response structures with validation tags
    - `service.go` - Business logic with audit logging
    - `handler.go` - Single handler (if module is used by one party only)
    - `handler_admin.go` - Admin handlers (if module is used by both admin and users)
    - `handler_frontend.go` - User-facing handlers (if module is used by both admin and users)
- `internal/audit/` - Audit and error logging services
- `internal/db/` - SQLC generated code (don't edit manually)
- `internal/middleware/` - HTTP middleware
- `internal/router/` - Route definitions
  - `router.go` - Main router setup
  - `admin.go` - Admin routes (protected by admin auth middleware)
  - `frontend.go` - User-facing routes (protected by user auth middleware)
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
