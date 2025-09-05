package handler

import (
	"net/http"
	"regexp"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/jwt"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	tokenService *jwt.TokenService
	logger       *logger.Logger
	config       *config.AuthConfig
}

// AuthRequest represents the authentication request
type AuthRequest struct {
	Email string `json:"email" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(tokenService *jwt.TokenService, logger *logger.Logger, config *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
		logger:       logger,
		config:       config,
	}
}

// Authenticate handles the authentication request
// @Summary Authenticate user
// @Description Authenticate a user and return a JWT token for API access
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Authentication request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth [post]
func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Validate email
	if req.Email == "" {
		h.logger.Error("Email is required", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}
	
	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		h.logger.Error("Invalid email format", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	// Generate token
	token, err := h.tokenService.GenerateToken(req.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Log successful authentication
	h.logger.Info("User authenticated", map[string]interface{}{
		"email": req.Email,
	})

	// Return token
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}
