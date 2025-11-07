package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/user/coc/internal/db"
	"github.com/user/coc/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	queries             *db.Queries
	jwtSecret           string
	bearerTokenDuration time.Duration
}

func NewService(queries *db.Queries, jwtSecret string, bearerTokenDuration time.Duration) *Service {
	return &Service{
		queries:             queries,
		jwtSecret:           jwtSecret,
		bearerTokenDuration: bearerTokenDuration,
	}
}

// CustomClaims defines the JWT claims structure
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (string, *db.User, error) {
	// Get user by email
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.Unauthorized("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.Unauthorized("invalid email or password")
	}

	// Generate JWT token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return "", nil, errors.Internal("failed to generate token", err)
	}

	return token, &user, nil
}

// GenerateToken creates a new JWT token for a user
func (s *Service) GenerateToken(user *db.User) (string, error) {
	// Token expires based on BEARER_TOKEN_DURATION from config
	expirationTime := time.Now().Add(s.bearerTokenDuration)

	// Convert UUID to string
	userID, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		return "", err
	}

	// Create claims
	claims := &CustomClaims{
		UserID:   userID.String(),
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*CustomClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, email, username, password, firstName, lastName string) (*db.User, error) {
	// Check if user with email already exists
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.AlreadyExists("user with this email already exists")
	}

	// Check if user with username already exists
	_, err = s.queries.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, errors.AlreadyExists("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Internal("failed to hash password", err)
	}

	// Create user
	params := db.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if firstName != "" {
		params.FirstName = pgtype.Text{String: firstName, Valid: true}
	}

	if lastName != "" {
		params.LastName = pgtype.Text{String: lastName, Valid: true}
	}

	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, errors.Internal("failed to create user", err)
	}

	return &user, nil
}
