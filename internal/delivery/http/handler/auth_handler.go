package handler

import (
	"net/http"
	"time"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/dex"
	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	dexClient *dex.Client
	logger    *logger.Logger
	config    *config.AuthConfig
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(dexClient *dex.Client, logger *logger.Logger, config *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		dexClient: dexClient,
		logger:    logger,
		config:    config,
	}
}

// Login redirects the user to the Dex login page
func (h *AuthHandler) Login(c *gin.Context) {
	// Generate a random state for CSRF protection
	state, err := h.dexClient.GenerateState()
	if err != nil {
		h.logger.Error("Failed to generate state", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate state",
		})
		return
	}

	// Store the state in a cookie
	c.SetCookie("auth_state", state, int(time.Hour.Seconds()), "/", "", false, true)

	// Redirect to the authorization URL
	authURL := h.dexClient.GetAuthURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// Callback handles the callback from the Dex server
func (h *AuthHandler) Callback(c *gin.Context) {
	// Get the state from the cookie
	state, err := c.Cookie("auth_state")
	if err != nil {
		h.logger.Error("Failed to get state cookie", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state",
		})
		return
	}

	// Verify the state
	if c.Query("state") != state {
		h.logger.Error("State mismatch", map[string]interface{}{
			"expected": state,
			"received": c.Query("state"),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "State mismatch",
		})
		return
	}

	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		h.logger.Error("Authorization code is missing", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code is missing",
		})
		return
	}

	// Exchange the code for tokens
	token, err := h.dexClient.Exchange(c.Request.Context(), code)
	if err != nil {
		h.logger.Error("Failed to exchange code for token", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange code for token",
		})
		return
	}

	// Get the ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		h.logger.Error("ID token is missing", nil)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ID token is missing",
		})
		return
	}

	// Verify the ID token and extract user information
	userInfo, err := h.dexClient.VerifyIDToken(c.Request.Context(), rawIDToken)
	if err != nil {
		h.logger.Error("Failed to verify ID token", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify ID token",
		})
		return
	}

	// Log the successful authentication
	h.logger.Info("User authenticated", map[string]interface{}{
		"email": userInfo.Email,
		"name":  userInfo.Name,
	})

	// Return the token and user information
	c.JSON(http.StatusOK, gin.H{
		"token":     rawIDToken,
		"tokenType": "Bearer",
		"expiresIn": userInfo.Expiry - time.Now().Unix(),
		"user": gin.H{
			"email":   userInfo.Email,
			"name":    userInfo.Name,
			"subject": userInfo.Subject,
		},
	})
}
