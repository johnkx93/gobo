# Admin Menu Tree - Implementation Examples

## Quick Start

### 1. Login as Admin and Get Menu

```bash
# Login
curl -X POST http://localhost:8080/api/admin/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123"
  }'

# Response
{
  "status": "success",
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "admin": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "admin@example.com",
      "username": "admin",
      "role": "super_admin",
      "is_active": true
    }
  }
}

# Get Menu (using the token)
curl http://localhost:8080/api/admin/v1/menu \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

# Response - Menu based on role
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
      }
      // ... more menu items
    ]
  }
}
```

### 2. Frontend Integration Example (React)

```typescript
// types/menu.ts
export interface MenuItem {
  id: string;
  label: string;
  icon?: string;
  path?: string;
  children?: MenuItem[];
  order: number;
}

// hooks/useAdminMenu.ts
import { useState, useEffect } from 'react';
import { MenuItem } from '../types/menu';

export const useAdminMenu = () => {
  const [menu, setMenu] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMenu = async () => {
      try {
        const token = localStorage.getItem('adminToken');
        const response = await fetch('/api/admin/v1/menu', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!response.ok) {
          throw new Error('Failed to fetch menu');
        }

        const { data } = await response.json();
        setMenu(data.menu);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchMenu();
  }, []);

  return { menu, loading, error };
};

// components/Sidebar.tsx
import { useAdminMenu } from '../hooks/useAdminMenu';
import { MenuItem } from '../types/menu';

const MenuItemComponent = ({ item }: { item: MenuItem }) => {
  if (item.children && item.children.length > 0) {
    return (
      <div className="menu-group">
        <div className="menu-header">
          {item.icon && <Icon name={item.icon} />}
          <span>{item.label}</span>
        </div>
        <div className="menu-children">
          {item.children.map(child => (
            <MenuItemComponent key={child.id} item={child} />
          ))}
        </div>
      </div>
    );
  }

  return (
    <Link to={item.path || '#'} className="menu-item">
      {item.icon && <Icon name={item.icon} />}
      <span>{item.label}</span>
    </Link>
  );
};

export const Sidebar = () => {
  const { menu, loading, error } = useAdminMenu();

  if (loading) return <div>Loading menu...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <nav className="sidebar">
      {menu.map(item => (
        <MenuItemComponent key={item.id} item={item} />
      ))}
    </nav>
  );
};
```

### 3. Testing Different Roles

```bash
# Test as Super Admin (sees everything)
curl http://localhost:8080/api/admin/v1/menu \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN"

# Test as Admin (limited access)
curl http://localhost:8080/api/admin/v1/menu \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Test as Moderator (read-only)
curl http://localhost:8080/api/admin/v1/menu \
  -H "Authorization: Bearer MODERATOR_TOKEN"
```

### 4. Permission-Protected Endpoint Examples

```bash
# Super admin can delete users (has users.delete permission)
curl -X DELETE http://localhost:8080/api/admin/v1/users/USER_ID \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN"
# Response: 200 OK

# Regular admin cannot delete users (lacks users.delete permission)
curl -X DELETE http://localhost:8080/api/admin/v1/users/USER_ID \
  -H "Authorization: Bearer ADMIN_TOKEN"
# Response: 403 Forbidden
{
  "status": "error",
  "message": "insufficient permissions for this action"
}

# Moderator cannot manage admins (lacks admins.manage permission)
curl http://localhost:8080/api/admin/v1/admins \
  -H "Authorization: Bearer MODERATOR_TOKEN"
# Response: 403 Forbidden
{
  "status": "error",
  "message": "insufficient permissions for this action"
}
```

## Role-Specific Menu Examples

### Super Admin Menu
```json
{
  "role": "super_admin",
  "menu": [
    {
      "id": "users",
      "label": "User Management",
      "icon": "users",
      "children": [
        { "id": "users-create", "label": "Create User", "path": "/admin/users/create" },
        { "id": "users-list", "label": "User List", "path": "/admin/users" }
      ]
    },
    {
      "id": "orders",
      "label": "Order Management",
      "icon": "shopping-cart",
      "children": [
        { "id": "orders-list", "label": "Order List", "path": "/admin/orders" },
        { "id": "orders-update", "label": "Update Orders", "path": "/admin/orders/bulk-update" }
      ]
    },
    {
      "id": "admins",
      "label": "Admin Management",
      "icon": "shield",
      "path": "/admin/admins"
    },
    {
      "id": "settings",
      "label": "Settings",
      "icon": "settings",
      "children": [
        { "id": "settings-general", "label": "General Settings", "path": "/admin/settings/general" },
        { "id": "settings-security", "label": "Security", "path": "/admin/settings/security" }
      ]
    }
  ]
}
```

### Regular Admin Menu
```json
{
  "role": "admin",
  "menu": [
    {
      "id": "users",
      "label": "User Management",
      "icon": "users",
      "children": [
        { "id": "users-create", "label": "Create User", "path": "/admin/users/create" },
        { "id": "users-list", "label": "User List", "path": "/admin/users" }
      ]
    },
    {
      "id": "orders",
      "label": "Order Management",
      "icon": "shopping-cart",
      "children": [
        { "id": "orders-list", "label": "Order List", "path": "/admin/orders" },
        { "id": "orders-update", "label": "Update Orders", "path": "/admin/orders/bulk-update" }
      ]
    },
    {
      "id": "settings",
      "label": "Settings",
      "icon": "settings",
      "children": [
        { "id": "settings-general", "label": "General Settings", "path": "/admin/settings/general" }
      ]
    }
  ]
}
```
Note: Admin Management and Security settings are hidden for regular admin.

### Moderator Menu
```json
{
  "role": "moderator",
  "menu": [
    {
      "id": "users",
      "label": "User Management",
      "icon": "users",
      "children": [
        { "id": "users-list", "label": "User List", "path": "/admin/users" }
      ]
    },
    {
      "id": "orders",
      "label": "Order Management",
      "icon": "shopping-cart",
      "children": [
        { "id": "orders-list", "label": "Order List", "path": "/admin/orders" }
      ]
    }
  ]
}
```
Note: Only read operations are available. No create/update/delete options.

## Adding Custom Menu Items

### Step 1: Update menu.go

```go
// Add to the allMenus slice in GetMenuForRole function
{
    ID:    "analytics",
    Label: "Analytics",
    Icon:  "chart-bar",
    Order: 5,
    Children: []MenuItem{
        {
            ID:         "analytics-dashboard",
            Label:      "Dashboard",
            Path:       "/admin/analytics/dashboard",
            Permission: "analytics.dashboard",
            Order:      1,
        },
        {
            ID:         "analytics-reports",
            Label:      "Reports",
            Path:       "/admin/analytics/reports",
            Permission: "analytics.reports",
            Order:      2,
        },
    },
}
```

### Step 2: Add Permissions

```go
// In getRolePermissions function
"super_admin": {
    // ... existing permissions
    "analytics.dashboard": true,
    "analytics.reports":   true,
},
"admin": {
    // ... existing permissions
    "analytics.dashboard": true,
},
```

### Step 3: Create Routes and Handlers

```go
// In router/admin.go
r.Route("/analytics", func(r chi.Router) {
    r.Use(adminAuthMiddleware)
    r.With(middleware.RequirePermission("analytics.dashboard")).
        Get("/dashboard", analyticsHandler.GetDashboard)
    r.With(middleware.RequirePermission("analytics.reports")).
        Get("/reports", analyticsHandler.GetReports)
})
```

## Best Practices

1. **Always fetch menu after login** - Menu changes based on role
2. **Cache menu in frontend** - Store in localStorage/state management
3. **Refresh on role change** - If admin role is updated, refetch menu
4. **Handle permission errors gracefully** - Show user-friendly messages
5. **Use optimistic UI** - Show menu immediately from cache, refresh in background
6. **Breadcrumbs from menu** - Use menu structure for auto-generated breadcrumbs
7. **Active state tracking** - Highlight current menu item based on route

## Security Notes

- ✅ Menu is filtered server-side based on role
- ✅ Each endpoint has permission middleware protection
- ✅ Cannot bypass permissions by calling API directly
- ✅ Tokens contain role information validated on each request
- ✅ Menu visibility and API access are synchronized
