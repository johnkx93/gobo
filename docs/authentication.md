# JWT Authentication Guide

This application uses JWT (JSON Web Tokens) for authentication with email/password login.

## Features

- **Email/Password Authentication**: Secure login with bcrypt password hashing
- **7-Day Token Expiration**: JWT tokens are valid for 7 days
- **Protected Routes**: User and Order endpoints require authentication
- **Automatic User Context**: Authenticated user information is available in request context

## Configuration

Add the following to your `.env` file:

```env
JWT_SECRET=your-secret-key-here
```

**Generate a secure secret:**
```bash
openssl rand -base64 32
```

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Register a New User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "status": true,
  "message": "registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2025-11-06T10:00:00Z",
      "updated_at": "2025-11-06T10:00:00Z"
    }
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response:**
```json
{
  "status": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2025-11-06T10:00:00Z",
      "updated_at": "2025-11-06T10:00:00Z"
    }
  }
}
```

### Protected Endpoints (Authentication Required)

All user and order endpoints require authentication. Include the JWT token in the `Authorization` header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Example: List Users
```http
GET /api/v1/users?limit=10&offset=0
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Example: Create Order
```http
POST /api/v1/orders
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "total_amount": 99.99,
  "notes": "Express delivery"
}
```

## Using the Authentication in Code

### Accessing Authenticated User

In your handlers, you can access the authenticated user from the request context:

```go
import "github.com/user/coc/internal/app/coc/auth"

func (h *Handler) MyHandler(w http.ResponseWriter, r *http.Request) {
    user := auth.UserFromContext(r)
    if user == nil {
        // User not authenticated (shouldn't happen on protected routes)
        return
    }
    
    // Use user.ID, user.Email, etc.
    fmt.Printf("Authenticated user: %s\n", user.Email)
}
```

### Getting User ID

```go
func (h *Handler) MyHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)
    // Use userID
}
```

## Security Features

1. **Password Hashing**: Passwords are hashed using bcrypt before storage
2. **No Password Exposure**: `password_hash` is never included in API responses
3. **Token Validation**: All protected routes validate JWT tokens
4. **User Verification**: Middleware verifies user still exists in database
5. **Secure Signing**: Tokens are signed with HS256 algorithm

## Token Structure

JWT tokens contain the following claims:

```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "username": "johndoe",
  "exp": 1699257600,  // Expiration time (7 days from issuance)
  "iat": 1698652800,  // Issued at
  "nbf": 1698652800   // Not before
}
```

## Error Responses

### Unauthorized (401)
```json
{
  "status": false,
  "message": "invalid email or password"
}
```

### Missing Token (401)
```json
{
  "status": false,
  "message": "missing authorization header"
}
```

### Expired Token (401)
```json
{
  "status": false,
  "message": "invalid or expired token"
}
```

### Conflict (409) - User Already Exists
```json
{
  "status": false,
  "message": "user with this email already exists"
}
```

## Testing with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Access Protected Endpoint
```bash
# Save token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

## Best Practices

1. **Store tokens securely** on the client side (e.g., httpOnly cookies or secure storage)
2. **Use HTTPS** in production to prevent token interception
3. **Rotate JWT_SECRET** periodically in production
4. **Implement token refresh** for better security (future enhancement)
5. **Set up rate limiting** on auth endpoints to prevent brute force attacks

## Future Enhancements

Consider implementing these features:

- [ ] Refresh token pattern for better security
- [ ] Token blacklist for logout functionality
- [ ] Password reset functionality
- [ ] Email verification
- [ ] Multi-factor authentication (MFA)
- [ ] Rate limiting on login attempts
- [ ] Account lockout after failed attempts
