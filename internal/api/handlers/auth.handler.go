package handlers

import (
	"errors"
	"net/http"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"     // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
	routes      *constants.Routes // <-- ADD THIS
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// @Summary      User Login
// @Description  Authenticates a user with email and password, returning a JWT and session details.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body ports.LoginInput true "User Login Credentials"
// @Success      200 {object} ports.LoginResponse
// @Failure      400 {object} map[string]string "Validation error"
// @Failure      401 {object} map[string]string "Invalid credentials or deactivated account"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input ports.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	response, err := h.authService.Login(c.Request.Context(), input, &ipAddress, &userAgent)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      User Logout
// @Description  Revokes the current user's session token.
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// --- USE CONSTANTS ---
	userIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	sessionIDVal, _ := c.Get(h.routes.ContextKeySessionID)

	userID, _ := userIDVal.(uint)
	sessionID, _ := sessionIDVal.(uint)

	if userID == 0 || sessionID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}

	_, err := h.authService.Logout(c.Request.Context(), sessionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Get Current User
// @Description  Retrieves the details of the currently authenticated user.
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} ports.UserResponse
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// --- USE CONSTANT ---
	userVal, exists := c.Get(h.routes.ContextKeyUser)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userVal.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	c.JSON(http.StatusOK, ports.MapUserToResponse(user))
}
