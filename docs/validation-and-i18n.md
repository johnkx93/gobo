# Validation and i18n Guide

This document explains how validation works in this project, including i18n support, custom validators, and custom messages.

## Overview

The project uses `go-playground/validator/v10` with universal-translator for multilingual validation error messages. The validator is initialized once in `cmd/api/main.go` and shared across all handlers.

**Supported Languages:**
- English (en)
- Malay (ms)
- Simplified Chinese (zh)

## How It Works

### 1. Centralized Validator Instance

Instead of creating `validator.New()` in each handler, we create one shared instance in `main.go`:

```go
// cmd/api/main.go
validator := validation.New()

// Pass to handlers
userHandler := user.NewHandler(userService, validator)
orderHandler := order.NewHandler(orderService, validator)
authHandler := auth.NewHandler(authService, validator)
```

**Benefits:**
- Single initialization (faster startup)
- Consistent validation rules across the app
- Custom validators registered once
- Translations loaded once

### 2. Using Validation Tags

Add validation rules to struct fields using tags:

```go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Password  string `json:"password" validate:"required,min=6"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
}
```

**Common Tags:**
- `required` - field must not be empty
- `email` - must be valid email format
- `min=N` - minimum length (strings) or value (numbers)
- `max=N` - maximum length or value
- `len=N` - exact length
- `oneof=a b c` - value must be one of the listed options
- `url` - must be valid URL
- `uuid` - must be valid UUID
- `omitempty` - skip validation if field is empty

### 3. Validating in Handlers

In your handler, call `Struct()` and use `TranslateErrors()` to get localized messages:

```go
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if err := h.validate.Struct(req); err != nil {
        errors := h.validate.TranslateErrors(r, err)
        h.respondJSON(w, http.StatusBadRequest, "validation failed", map[string]interface{}{"errors": errors})
        return
    }

    // Proceed with business logic...
}
```

### 4. i18n Support (Multilingual Error Messages)

The validator automatically detects the client's preferred language from the `Accept-Language` HTTP header.

**Example Requests:**

```bash
# English (default)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"email": "", "username": "ab", "password": "123"}'

# Response:
{
  "status": false,
  "message": "validation failed",
  "data": {
    "errors": {
      "email": "email is required and cannot be empty",
      "username": "username must be at least 3 characters long",
      "password": "password must be at least 6 characters long"
    }
  }
}
```

```bash
# Malay
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ms" \
  -d '{"email": "", "username": "ab", "password": "123"}'

# Response:
{
  "status": false,
  "message": "validation failed",
  "data": {
    "errors": {
      "email": "email adalah wajib",
      "username": "username mesti sekurang-kurangnya 3 aksara",
      "password": "password mesti sekurang-kurangnya 6 aksara"
    }
  }
}
```

```bash
# Simplified Chinese
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh" \
  -d '{"email": "", "username": "ab", "password": "123"}'

# Response (Chinese error messages)
{
  "status": false,
  "message": "validation failed",
  "data": {
    "errors": {
      "email": "email为必填字段",
      "username": "username长度必须至少为3个字符",
      "password": "password长度必须至少为6个字符"
    }
  }
}
```

## Custom Validators

The project includes two example custom validators in `internal/validation/validator.go`:

### 1. `notbadword` - Content Filter

Checks if a string contains inappropriate words.

**Usage:**
```go
type CommentRequest struct {
    Content string `json:"content" validate:"required,notbadword"`
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"content": "This is spam content"}'

# Response:
{
  "errors": {
    "content": "content contains inappropriate content"
  }
}
```

### 2. `strongpassword` - Password Strength

Validates password has:
- At least 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit

**Usage:**
```go
type SetPasswordRequest struct {
    Password string `json:"password" validate:"required,strongpassword"`
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/password \
  -H "Content-Type: application/json" \
  -d '{"password": "weak"}'

# Response:
{
  "errors": {
    "password": "password must be at least 8 characters and contain uppercase, lowercase, and digit"
  }
}
```

### Adding Your Own Custom Validator

Edit `internal/validation/validator.go` and add to `registerCustomValidators()`:

```go
func (v *Validator) registerCustomValidators() {
    // Your custom validator
    v.validate.RegisterValidation("yourvalidator", func(fl validator.FieldLevel) bool {
        value := fl.Field().String()
        // Your validation logic
        return true // or false if invalid
    })
}
```

Then register translations for it in `registerCustomMessages()`:

```go
// English
enTrans, _ := v.uni.GetTranslator("en")
v.validate.RegisterTranslation("yourvalidator", enTrans,
    func(ut ut.Translator) error {
        return ut.Add("yourvalidator", "{0} has custom validation error", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("yourvalidator", fe.Field())
        return t
    },
)

// Chinese
zhTrans, _ := v.uni.GetTranslator("zh")
v.validate.RegisterTranslation("yourvalidator", zhTrans,
    func(ut ut.Translator) error {
        return ut.Add("yourvalidator", "{0}自定义验证错误", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("yourvalidator", fe.Field())
        return t
    },
)
```

## Custom Error Messages

You can override the default message for any validation tag. The project overrides messages for `required`, `email`, and `min` tags.

**Example from `internal/validation/validator.go`:**

```go
func (v *Validator) registerCustomMessages() {
    enTrans, _ := v.uni.GetTranslator("en")
    
    // Custom message for "required" tag
    v.validate.RegisterTranslation("required", enTrans,
        func(ut ut.Translator) error {
            return ut.Add("required", "{0} is required and cannot be empty", true)
        },
        func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("required", fe.Field())
            return t
        },
    )
    
    // Custom message for "email" tag
    v.validate.RegisterTranslation("email", enTrans,
        func(ut ut.Translator) error {
            return ut.Add("email", "{0} must be a valid email address", true)
        },
        func(ut ut.Translator, fe validator.FieldError) string {
            t, _ := ut.T("email", fe.Field())
            return t
        },
    )
}
```

### Message Placeholders

- `{0}` - Field name (from JSON tag)
- `{1}` - Parameter (e.g., min value, max value)

**Example with parameter:**
```go
v.validate.RegisterTranslation("min", enTrans,
    func(ut ut.Translator) error {
        return ut.Add("min", "{0} must be at least {1} characters long", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("min", fe.Field(), fe.Param())
        return t
    },
)
```

## JSON Field Names in Errors

The validator is configured to use JSON field names (not struct field names) in error messages:

```go
validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if name == "" || name == "-" {
        return fld.Name
    }
    return name
})
```

**Example:**

```go
type User struct {
    Email string `json:"email" validate:"required"`
}
```

Error message will say `"email is required"` (not `"Email is required"`).

## Testing Validation

### Test with curl

```bash
# Test required field
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{}'

# Test email format
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid", "username": "test", "password": "password123"}'

# Test minimum length
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "username": "ab", "password": "pass"}'

# Test with Malay language
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ms" \
  -d '{"email": "", "username": "", "password": ""}'

# Test with Chinese language
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh-CN" \
  -d '{"email": "", "username": "", "password": ""}'
```

## Advanced Usage

### Struct-Level Validation

For validation that depends on multiple fields:

```go
func (v *Validator) registerCustomValidators() {
    v.validate.RegisterStructValidation(func(sl validator.StructLevel) {
        req := sl.Current().Interface().(RegisterRequest)
        
        // Example: username cannot equal password
        if req.Username == req.Password {
            sl.ReportError(req.Password, "password", "Password", "nefield", "username")
        }
    }, RegisterRequest{})
}
```

### Conditional Validation

Use tags like:
- `required_if=Field Value` - required if Field equals Value
- `required_with=Field` - required if Field is present
- `required_without=Field` - required if Field is absent

```go
type UpdateUserRequest struct {
    Password    string `json:"password,omitempty"`
    NewPassword string `json:"new_password" validate:"required_with=Password,min=8"`
}
```

### Cross-Field Validation

```go
type ChangePasswordRequest struct {
    Password        string `json:"password" validate:"required"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
```

## Best Practices

1. **Use `omitempty` in validate tag for optional fields:**
   ```go
   Bio string `json:"bio" validate:"omitempty,max=500"`
   ```

2. **Normalize before validating:**
   ```go
   req.Email = strings.TrimSpace(strings.ToLower(req.Email))
   if err := h.validate.Struct(req); err != nil { ... }
   ```

3. **Don't put DB checks in validators:**
   - Validators should be stateless and fast
   - Do DB uniqueness checks in the service layer

4. **Keep custom validators simple:**
   - Focus on format/content validation
   - Avoid I/O operations

5. **Provide clear error messages:**
   - Use descriptive field names in JSON tags
   - Write user-friendly custom messages

## References

- [go-playground/validator documentation](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [universal-translator](https://pkg.go.dev/github.com/go-playground/universal-translator)
- [Validation tags reference](https://github.com/go-playground/validator#baked-in-validations)
