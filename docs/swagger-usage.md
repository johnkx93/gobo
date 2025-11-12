# Swagger Documentation Guide

## Overview

This project uses Swaggo to automatically generate OpenAPI/Swagger documentation from code annotations. The Swagger UI provides an interactive interface to explore and test API endpoints.

## Accessing Swagger UI

Once the application is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Generating Documentation

### Automatic Generation

Swagger documentation is generated from code annotations. After modifying handlers or DTOs, regenerate the docs:

```bash
make swagger
```

This command:

1. Parses annotations in your code
2. Generates `docs/swagger/` directory with:
   - `docs.go` - Go documentation file
   - `swagger.json` - OpenAPI JSON spec
   - `swagger.yaml` - OpenAPI YAML spec

### Manual Generation

If `make swagger` doesn't work, use the full path:

```bash
$HOME/go/bin/swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal
```

## Annotating Code

### 1. General API Information

Located in `cmd/api/main.go`:

```go
// @title           BO API
// @version         1.0
// @description     Backend API for business operations with admin and frontend endpoints
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

### 2. Handler Annotations

Add annotations above each handler function:

```go
// GetMe returns the current user's profile
// @Summary      Get current user profile
// @Description  Retrieve the authenticated user's profile information
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Success      200 {object} response.JSONResponse{data=UserResponse} "User profile retrieved"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/users/me [get]
func (h *FrontendHandler) GetMe(w http.ResponseWriter, r *http.Request) {
    // Handler implementation...
}
```

### 3. DTO Examples

Add `example` tags to DTO fields for better documentation:

```go
type LoginRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50" example:"admin"`
    Password string `json:"password" validate:"required,min=6" example:"password123"`
}

type UserResponse struct {
    ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Email     string `json:"email" example:"john.doe@example.com"`
    Username  string `json:"username" example:"johndoe"`
    CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
}
```

## Annotation Tags Reference

### Common Tags

| Tag            | Description           | Example                                         |
| -------------- | --------------------- | ----------------------------------------------- |
| `@Summary`     | Short description     | `@Summary Get user profile`                     |
| `@Description` | Detailed description  | `@Description Retrieve authenticated user data` |
| `@Tags`        | Group endpoints       | `@Tags User Profile`                            |
| `@Accept`      | Request content type  | `@Accept json`                                  |
| `@Produce`     | Response content type | `@Produce json`                                 |
| `@Param`       | Request parameter     | `@Param id path string true "User ID"`          |
| `@Success`     | Success response      | `@Success 200 {object} UserResponse`            |
| `@Failure`     | Error response        | `@Failure 404 {object} ErrorResponse`           |
| `@Security`    | Authentication method | `@Security BearerAuth`                          |
| `@Router`      | Endpoint route        | `@Router /api/v1/users/{id} [get]`              |

### Parameter Types

```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter
// @Param page query int false "Page number"

// Header parameter
// @Param Authorization header string true "Bearer token"

// Body parameter
// @Param request body CreateUserRequest true "User data"
```

### Response Types

```go
// Simple response
// @Success 200 {object} UserResponse

// Response with array
// @Success 200 {object} response.JSONResponse{data=[]UserResponse}

// Response with nested data
// @Success 200 {object} response.JSONResponse{data=LoginResponse}

// Multiple success codes
// @Success 200 {object} UserResponse "OK"
// @Success 201 {object} UserResponse "Created"
```

## Workflow

### Adding a New Endpoint

1. **Create handler function** in `internal/app/*/handler.go`
2. **Add Swagger annotations** above the function
3. **Update DTOs** with `example` tags if needed
4. **Regenerate docs**: `make swagger`
5. **Test endpoint** in Swagger UI at `http://localhost:8080/swagger/index.html`

### Example: Complete Annotation

```go
// CreateAddress handles POST /api/v1/addresses
// @Summary      Create user address
// @Description  Create a new address for the authenticated user
// @Tags         User Addresses
// @Accept       json
// @Produce      json
// @Param        request body UserCreateAddressRequest true "Address data"
// @Success      201 {object} response.JSONResponse{data=[]AddressResponse} "Address created successfully"
// @Failure      400 {object} response.JSONResponse "Invalid request"
// @Failure      401 {object} response.JSONResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/addresses [post]
func (h *FrontendHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
    // Implementation...
}
```

## Testing Endpoints

### Using Swagger UI

1. Navigate to `http://localhost:8080/swagger/index.html`
2. Browse endpoints by tags (Admin Authentication, User Profile, etc.)
3. Click on an endpoint to expand it
4. Click **"Try it out"**
5. Fill in parameters (if required)
6. For authenticated endpoints:
   - Click **"Authorize"** button at the top
   - Enter: `Bearer YOUR_JWT_TOKEN`
   - Click **"Authorize"** then **"Close"**
7. Click **"Execute"**
8. View response below

### Example: Testing Login

1. Find **Admin Authentication** → **POST /api/admin/v1/auth/login**
2. Click **"Try it out"**
3. Edit request body:
   ```json
   {
     "username": "admin",
     "password": "your_password"
   }
   ```
4. Click **"Execute"**
5. Copy the `token` from the response
6. Click **"Authorize"** and paste: `Bearer <token>`
7. Now you can test authenticated endpoints

## Troubleshooting

### Swagger UI shows old endpoints

```bash
make swagger  # Regenerate documentation
# Then restart the server
make run
```

### Cannot find response type error

Ensure you're using the correct type names:

- ✅ `response.JSONResponse` (not `response.Response`)
- Check imports in handlers

### Swag command not found

Install the Swag CLI tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Or use the full path in Makefile:

```bash
$(HOME)/go/bin/swag init -g cmd/api/main.go -o docs/swagger
```

### Build errors after adding Swagger

Check that imports are correct:

```go
import (
    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/user/coc/docs/swagger" // Blank import for side effects
)
```

## Best Practices

1. **Always add examples** to DTO fields for better documentation
2. **Group related endpoints** using consistent `@Tags`
3. **Document all error codes** with appropriate `@Failure` annotations
4. **Keep descriptions concise** but informative
5. **Regenerate docs** after every handler/DTO change
6. **Test in Swagger UI** before committing changes
7. **Use semantic versioning** in `@version` tag

## Advanced Features

### Multiple Response Schemas

```go
// @Success 200 {object} response.JSONResponse{data=UserResponse} "User found"
// @Success 200 {object} response.JSONResponse{data=[]UserResponse} "User list"
// @Failure 404 {object} response.JSONResponse "User not found"
// @Failure 500 {object} response.JSONResponse "Internal server error"
```

### Custom Headers

```go
// @Param X-Request-ID header string false "Request ID for tracking"
// @Success 200 {object} UserResponse
// @Header 200 {string} X-Request-ID "Request tracking ID"
```

### File Uploads

```go
// @Param file formData file true "File to upload"
// @Success 200 {object} response.JSONResponse{data=FileResponse}
```

## Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
