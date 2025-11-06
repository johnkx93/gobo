# Quick Validation Examples

## Testing the New Validation System

### 1. Test i18n with Different Languages

```bash
# English (default)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d '{"email": "notanemail", "username": "ab", "password": "123"}'

# Malay
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ms" \
  -d '{"email": "notanemail", "username": "ab", "password": "123"}'

# Simplified Chinese
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh" \
  -d '{"email": "notanemail", "username": "ab", "password": "123"}'
```

### 2. Test Custom Validator: notbadword

To test the `notbadword` custom validator, add it to a DTO. For example, update `RegisterRequest`:

```go
// internal/auth/dto.go
type RegisterRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Username  string `json:"username" validate:"required,min=3,max=50,notbadword"`
    Password  string `json:"password" validate:"required,min=6"`
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
}
```

Then test:

```bash
# This will fail validation
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "badword",
    "password": "password123"
  }'

# Response:
{
  "status": false,
  "message": "validation failed",
  "data": {
    "errors": {
      "username": "username contains inappropriate content"
    }
  }
}
```

### 3. Test Custom Validator: strongpassword

Add to a password field:

```go
// Example in dto.go
type SetPasswordRequest struct {
    Password string `json:"password" validate:"required,strongpassword"`
}
```

Test:

```bash
# Weak password (fails)
curl -X POST http://localhost:8080/api/v1/users/password \
  -H "Content-Type: application/json" \
  -d '{"password": "weak"}'

# Response:
{
  "errors": {
    "password": "password must be at least 8 characters and contain uppercase, lowercase, and digit"
  }
}

# Strong password (passes)
curl -X POST http://localhost:8080/api/v1/users/password \
  -H "Content-Type: application/json" \
  -d '{"password": "Strong123"}'
```

### 4. Adding a Custom Validator (Example)

Let's say you want to add a validator for phone numbers:

**Step 1:** Add to `internal/validation/validator.go` in `registerCustomValidators()`:

```go
// Custom validator for Malaysian phone numbers
v.validate.RegisterValidation("myphone", func(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    // Simple Malaysian phone format: +60xxxxxxxxx or 01xxxxxxxx
    if strings.HasPrefix(phone, "+60") && len(phone) >= 12 {
        return true
    }
    if strings.HasPrefix(phone, "01") && len(phone) >= 10 {
        return true
    }
    return false
})
```

**Step 2:** Add translations in `registerCustomMessages()`:

```go
// English
enTrans, _ := v.uni.GetTranslator("en")
v.validate.RegisterTranslation("myphone", enTrans,
    func(ut ut.Translator) error {
        return ut.Add("myphone", "{0} must be a valid Malaysian phone number", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("myphone", fe.Field())
        return t
    },
)

// Chinese
zhTrans, _ := v.uni.GetTranslator("zh")
v.validate.RegisterTranslation("myphone", zhTrans,
    func(ut ut.Translator) error {
        return ut.Add("myphone", "{0}必须是有效的马来西亚电话号码", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("myphone", fe.Field())
        return t
    },
)
```

**Step 3:** Use in your DTO:

```go
type UpdateProfileRequest struct {
    Phone string `json:"phone" validate:"omitempty,myphone"`
}
```

**Step 4:** Test:

```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ms" \
  -d '{"phone": "123456"}'

# Response (in Malay):
{
  "errors": {
    "phone": "phone mesti nombor telefon Malaysia yang sah"
  }
}
```

### 5. Changing Message for Built-in Validators

To change the message for a built-in validator like `max`, add to `registerCustomMessages()`:

```go
// Custom message for "max" tag
enTrans, _ := v.uni.GetTranslator("en")
v.validate.RegisterTranslation("max", enTrans,
    func(ut ut.Translator) error {
        return ut.Add("max", "{0} cannot exceed {1} characters", true)
    },
    func(ut ut.Translator, fe validator.FieldError) string {
        t, _ := ut.T("max", fe.Field(), fe.Param())
        return t
    },
)
```

Then restart your app and the new message will be used.

## Testing the Full Flow

### Start the server with debug logging:

```bash
LOG_LEVEL=debug make run
```

### Test validation with different scenarios:

```bash
# 1. All fields empty (test required)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{}'

# 2. Invalid email format
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid", "username": "testuser", "password": "password123"}'

# 3. Username too short (min=3)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "username": "ab", "password": "password123"}'

# 4. Password too short (min=6)
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "username": "testuser", "password": "123"}'

# 5. Valid request (should pass)
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

### Test language switching:

```bash
# Same invalid request in 3 languages
INVALID_DATA='{"email": "", "username": "ab", "password": "123"}'

echo "=== English ==="
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: en" \
  -d "$INVALID_DATA"

echo "\n\n=== Malay ==="
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: ms" \
  -d "$INVALID_DATA"

echo "\n\n=== Chinese ==="
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Accept-Language: zh" \
  -d "$INVALID_DATA"
```

## Summary of Changes

✅ **Centralized validator** - Created once in main.go, shared across all handlers
✅ **i18n support** - Automatic language detection from Accept-Language header
✅ **3 languages** - English, Malay, Simplified Chinese
✅ **Custom validators** - `notbadword` and `strongpassword` examples
✅ **Custom messages** - Override default messages for any tag
✅ **JSON field names** - Errors use JSON field names (not struct field names)
✅ **Structured errors** - Returns map[string]string for easy client parsing

## Files Modified

1. **Created:**
   - `internal/validation/validator.go` - Main validation package
   - `docs/validation-and-i18n.md` - Comprehensive documentation
   - `docs/validation-examples.md` - This file (quick examples)

2. **Modified:**
   - `cmd/api/main.go` - Initialize validator once
   - `internal/auth/handler.go` - Use shared validator + translations
   - `internal/app/user/handler.go` - Use shared validator + translations
   - `internal/app/order/handler.go` - Use shared validator + translations

3. **Dependencies Added:**
   - `github.com/go-playground/locales` - Locale definitions
   - `github.com/go-playground/universal-translator` - i18n support
   - Translation packages for en, ms, zh
