package admin_auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles admin authentication
type AuthService struct {
	queries   *db.Queries
	jwtSecret string
}

func NewAuthService(queries *db.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

// AdminClaims defines the JWT claims structure for admins
type AdminClaims struct {
	AdminID  string `json:"admin_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates an admin and returns a JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *db.Admin, error) {
	// Get admin by email
	admin, err := s.queries.GetAdminByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.Unauthorized("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.Unauthorized("invalid email or password")
	}

	// Check if admin is active
	if !admin.IsActive {
		return "", nil, errors.Unauthorized("admin account is disabled")
	}

	// Generate JWT token
	token, err := s.GenerateToken(&admin)
	if err != nil {
		return "", nil, errors.Internal("failed to generate token", err)
	}

	return token, &admin, nil
}

// GenerateToken creates a new JWT token for an admin
func (s *AuthService) GenerateToken(admin *db.Admin) (string, error) {
	// Token expires in 24 hours (typical for admin sessions)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Convert UUID to string
	adminID, err := uuid.FromBytes(admin.ID.Bytes[:])
	if err != nil {
		return "", err
	}

	// Create claims with admin role
	claims := &AdminClaims{
		AdminID:  adminID.String(),
		Email:    admin.Email,
		Username: admin.Username,
		Role:     admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "admin",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret (can use different secret for admin if desired)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the admin claims
func (s *AuthService) ValidateToken(tokenString string) (*AdminClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify this is an admin token
	if claims.Subject != "admin" {
		return nil, fmt.Errorf("not an admin token")
	}

	return claims, nil
}

// GetAdminByID retrieves an admin by ID (used by middleware)
func (s *AuthService) GetAdminByID(ctx context.Context, idStr string) (*AdminResponse, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.Validation("invalid admin ID format")
	}

	uuidBytes := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	admin, err := s.queries.GetAdminByID(ctx, uuidBytes)
	if err != nil {
		return nil, errors.NotFound("admin not found")
	}

	return toAdminResponse(&admin), nil
}

// Helper to convert db.Admin to AdminResponse
func toAdminResponse(admin *db.Admin) *AdminResponse {
	adminID, _ := uuid.FromBytes(admin.ID.Bytes[:])

	return &AdminResponse{
		ID:       adminID.String(),
		Email:    admin.Email,
		Username: admin.Username,
		Role:     admin.Role,
		IsActive: admin.IsActive,
	}
}

// CreateAdmin creates a new admin account (only super_admin can do this)
func (s *AuthService) CreateAdmin(ctx context.Context, email, username, password, firstName, lastName, role string) (*db.Admin, error) {
	// Check if admin with email already exists
	_, err := s.queries.GetAdminByEmail(ctx, email)
	if err == nil {
		return nil, errors.AlreadyExists("admin with this email already exists")
	}

	// Check if admin with username already exists
	_, err = s.queries.GetAdminByUsername(ctx, username)
	if err == nil {
		return nil, errors.AlreadyExists("admin with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Internal("failed to hash password", err)
	}

	// Create admin
	params := db.CreateAdminParams{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
	}

	if firstName != "" {
		params.FirstName = pgtype.Text{String: firstName, Valid: true}
	}

	if lastName != "" {
		params.LastName = pgtype.Text{String: lastName, Valid: true}
	}

	admin, err := s.queries.CreateAdmin(ctx, params)
	if err != nil {
		return nil, errors.Internal("failed to create admin", err)
	}

	return &admin, nil
}
