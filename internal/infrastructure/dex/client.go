package dex

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Client represents a Dex OIDC client
type Client struct {
	Provider     *oidc.Provider
	OAuth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	Logger       *logger.Logger
}

// UserInfo represents the user information extracted from the ID token
type UserInfo struct {
	Email    string `json:"email"`
	Subject  string `json:"sub"`
	Name     string `json:"name"`
	IssuedAt int64  `json:"iat"`
	Expiry   int64  `json:"exp"`
}

// NewClient creates a new Dex OIDC client
func NewClient(cfg *config.AuthConfig, log *logger.Logger) (*Client, error) {
	ctx := context.Background()

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, cfg.DexURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Configure OAuth2
	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &Client{
		Provider:     provider,
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
		Logger:       log,
	}, nil
}

// GenerateState generates a random state string for CSRF protection
func (c *Client) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GetAuthURL returns the authorization URL for the OIDC flow
func (c *Client) GetAuthURL(state string) string {
	return c.OAuth2Config.AuthCodeURL(state)
}

// Exchange exchanges the authorization code for tokens
func (c *Client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.OAuth2Config.Exchange(ctx, code)
}

// VerifyIDToken verifies the ID token and extracts user information
func (c *Client) VerifyIDToken(ctx context.Context, rawIDToken string) (*UserInfo, error) {
	idToken, err := c.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	var claims struct {
		Email   string `json:"email"`
		Subject string `json:"sub"`
		Name    string `json:"name"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	if claims.Email == "" {
		return nil, errors.New("email claim is missing from ID token")
	}

	return &UserInfo{
		Email:    claims.Email,
		Subject:  claims.Subject,
		Name:     claims.Name,
		IssuedAt: idToken.IssuedAt.Unix(),
		Expiry:   idToken.Expiry.Unix(),
	}, nil
}

// VerifyToken verifies a token from the Authorization header
func (c *Client) VerifyToken(ctx context.Context, r *http.Request) (*UserInfo, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	// Check if the header starts with "Bearer "
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return nil, errors.New("authorization header format must be 'Bearer {token}'")
	}

	// Extract the token
	rawIDToken := authHeader[len(prefix):]
	return c.VerifyIDToken(ctx, rawIDToken)
}

// IsSuperAdmin checks if the user is a super admin based on email
func (c *Client) IsSuperAdmin(email string, superAdminEmail string) bool {
	return email == superAdminEmail
}

// TokenValid checks if a token is still valid
func (c *Client) TokenValid(userInfo *UserInfo) bool {
	return time.Now().Unix() < userInfo.Expiry
}
