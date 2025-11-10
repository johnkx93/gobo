# Admin Panel Menu System

## Overview

The admin panel uses a role-based menu system that dynamically generates the navigation menu based on the authenticated admin's role and permissions.

## Architecture

### Components

1. **Menu Structure** (`internal/app/admin_auth/menu.go`)
   - Defines the complete menu hierarchy
   - Filters menus based on admin roles
   - Permission-based access control

2. **Menu Handler** (`internal/app/admin_auth/menu_handler.go`)
   - Provides API endpoint to fetch menu
   - Returns role-specific menu structure

3. **API Endpoint**: `GET /api/admin/v1/menu` (Protected)

## Menu Structure

```json
{
  "status": "success",
  "message": "menu retrieved successfully",
  "data": {
    "role": "super_admin",
    "menu": [
      {
        "id": "users",
        "label": "User Management",
        "icon": "users",
        "order": 1,
        "children": [
          {
            "id": "users-create",
            "label": "Create User",
            "path": "/admin/users/create",
            "order": 1
          },
          {
            "id": "users-list",
            "label": "User List",
            "path": "/admin/users",
            "order": 2
          }
        ]
      },
      {
        "id": "orders",
        "label": "Order Management",
        "icon": "shopping-cart",
        "order": 2,
        "children": [
          {
            "id": "orders-list",
            "label": "Order List",
            "path": "/admin/orders",
            "order": 1
          },
          {
            "id": "orders-update",
            "label": "Update Orders",
            "path": "/admin/orders/bulk-update",
            "order": 2
          }
        ]
      },
      {
        "id": "admins",
        "label": "Admin Management",
        "icon": "shield",
        "path": "/admin/admins",
        "order": 3
      },
      {
        "id": "settings",
        "label": "Settings",
        "icon": "settings",
        "order": 4,
        "children": [
          {
            "id": "settings-general",
            "label": "General Settings",
            "path": "/admin/settings/general",
            "order": 1
          },
          {
            "id": "settings-security",
            "label": "Security",
            "path": "/admin/settings/security",
            "order": 2
          }
        ]
      }
    ]
  }
}
```

## Role Permissions

### Super Admin
- Full access to all menus and features
- Permissions:
  - `users.*` (create, read, update, delete)
  - `orders.*` (read, update, delete)
  - `admins.manage`
  - `settings.*` (general, security)

### Admin
- Limited access to user and order management
- Permissions:
  - `users.create`, `users.read`, `users.update`
  - `orders.read`, `orders.update`
  - `settings.general`

### Moderator
- Read-only access
- Permissions:
  - `users.read`
  - `orders.read`

## Usage

### Frontend Integration

```javascript
// Fetch menu on admin login or app initialization
const fetchAdminMenu = async () => {
  const response = await fetch('/api/admin/v1/menu', {
    headers: {
      'Authorization': `Bearer ${adminToken}`
    }
  });
  
  const { data } = await response.json();
  return data.menu;
};

// Render menu dynamically
const renderMenu = (menuItems) => {
  return menuItems.map(item => ({
    ...item,
    // Frontend can add icons, routing logic, etc.
  }));
};
```

### Adding New Menu Items

1. **Add to menu structure** (`menu.go`):
```go
{
    ID:    "reports",
    Label: "Reports",
    Icon:  "chart",
    Order: 5,
    Children: []MenuItem{
        {
            ID:         "reports-sales",
            Label:      "Sales Report",
            Path:       "/admin/reports/sales",
            Permission: "reports.sales",
            Order:      1,
        },
    },
}
```

2. **Add permissions** to role definitions:
```go
"super_admin": {
    // ... existing permissions
    "reports.sales": true,
}
```

3. **Create the corresponding API endpoints** and handlers

## Advanced Features

### Database-Driven Menus (Future Enhancement)

For more dynamic menu management, consider storing menus in the database:

```sql
CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES menu_items(id),
    label VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    path VARCHAR(255),
    permission VARCHAR(100),
    order_index INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role VARCHAR(50) NOT NULL,
    permission VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(role, permission)
);
```

### Permission Middleware

Protect specific endpoints with permission checks:

```go
// In middleware/permissions.go
func RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            role := r.Context().Value(AdminRoleContextKey).(string)
            permissions := admin_auth.GetRolePermissions(role)
            
            if !permissions[permission] {
                response.Error(w, http.StatusForbidden, "insufficient permissions")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Usage in router
r.With(middleware.RequirePermission("users.delete")).Delete("/{id}", userAdminHandler.DeleteUser)
```

## Best Practices

1. **Keep menu structure simple**: Don't nest too deeply (max 2-3 levels)
2. **Use consistent naming**: Follow `resource.action` pattern for permissions
3. **Order matters**: Use the `order` field for consistent menu ordering
4. **Icon consistency**: Use a standard icon library (e.g., Heroicons, FontAwesome)
5. **Lazy loading**: For large menus, consider lazy loading sub-menus
6. **Cache menu**: Frontend should cache the menu and only refresh on role change
7. **Breadcrumbs**: Use menu structure to generate breadcrumbs automatically

## Testing

```bash
# Login as super_admin
curl -X POST http://localhost:8080/api/admin/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'

# Get menu (use token from login response)
curl http://localhost:8080/api/admin/v1/menu \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Alternative Approaches

### Option 2: Frontend-Defined Menu (Simpler)

If your frontend framework has good route guards, you can define the menu structure in the frontend and rely on backend API errors for permission enforcement:

**Pros:**
- Faster development
- No backend changes needed for menu updates
- Frontend has full control over UX

**Cons:**
- Security relies on API-level checks
- Menu items might show but fail when clicked
- Duplicated permission logic

### Option 3: Hybrid Approach

Define base menu in frontend, fetch visibility/permissions from backend:

```javascript
const menuConfig = [/* frontend menu */];
const { permissions } = await fetch('/api/admin/v1/permissions');
const visibleMenu = filterMenuByPermissions(menuConfig, permissions);
```

## Recommendation

For your current architecture, **Option 1 (Backend-Driven)** is recommended because:
- ✅ Single source of truth for permissions
- ✅ Role-based access control is already implemented
- ✅ Easy to extend with database-driven approach later
- ✅ Consistent with your audit logging and security model
- ✅ Prevents showing unauthorized options to users
