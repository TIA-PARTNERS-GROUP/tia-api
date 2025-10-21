package handlers

import (
	"errors"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
	routes      *constants.Routes
}

func NewAuthHandler(authService *services.AuthService, routes *constants.Routes) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
		routes:      routes,
	}
}
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
func (h *AuthHandler) Logout(c *gin.Context) {

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
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {

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
