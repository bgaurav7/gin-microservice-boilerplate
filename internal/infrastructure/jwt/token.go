package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Subject returns the subject from the registered claims
func (c *Claims) Subject() string {
	return c.RegisteredClaims.Subject
}

// TokenService handles JWT token operations
type TokenService struct {
	config *config.AuthConfig
}

// NewTokenService creates a new token service
func NewTokenService(config *config.AuthConfig) *TokenService {
	return &TokenService{
		config: config,
	}
}

// GenerateToken generates a new JWT token for the given email
func (s *TokenService) GenerateToken(email string) (string, error) {
	if email == "" {
		return "", errors.New("email cannot be empty")
	}

	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(s.config.JWTExpiryHours) * time.Hour)

	// Create claims
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	return claims, nil
}

// IsSuperAdmin checks if the user is a super admin based on email
func (s *TokenService) IsSuperAdmin(email string) bool {
	return email == s.config.SuperAdminEmail
}
