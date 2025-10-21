package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserConfigHandler struct {
	userConfigService *services.UserConfigService
	validate          *validator.Validate
	routes            *constants.Routes
}

func NewUserConfigHandler(userConfigService *services.UserConfigService, routes *constants.Routes) *UserConfigHandler {
	return &UserConfigHandler{
		userConfigService: userConfigService,
		validate:          validator.New(),
		routes:            routes,
	}
}

func (h *UserConfigHandler) getAuthUserID(c *gin.Context) (uint, error) {
	authUserIDVal, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		return 0, errors.New("invalid authentication context")
	}
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		return 0, errors.New("invalid authentication context")
	}
	return authUserID, nil
}

func (h *UserConfigHandler) checkUserOwnership(c *gin.Context) (uint, uint, *ports.ApiError) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		return 0, 0, ports.ErrInvalidToken
	}

	targetUserIDStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return 0, 0, ports.ErrDatabase
	}

	if authUserID != uint(targetUserID) {
		return 0, 0, ports.ErrForbidden
	}

	return authUserID, uint(targetUserID), nil
}

// @Summary Set or Update User Configuration
// @Description Creates a new configuration entry for a user, or updates an existing one for the given config_type. Enforces self-management.
// @Tags users, config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param config body ports.SetUserConfigInput true "Configuration Data"
// @Success 200 {object} ports.UserConfigResponse "Configuration successfully saved/updated"
// @Failure 400 {object} gin.H "Invalid request body or validation error"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden: Cannot modify another user's config"
// @Failure 500 {object} gin.H "Database error"
// @Router /users/{id}/config [put]
func (h *UserConfigHandler) SetUserConfig(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	var input ports.SetUserConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.userConfigService.SetUserConfig(c.Request.Context(), targetUserID, input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, ports.MapUserConfigToResponse(config))
}

// @Summary Get User Configuration by Type
// @Description Retrieves a specific configuration entry for a user. Enforces self-management.
// @Tags users, config
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param configType path string true "Configuration Type (e.g., user_preferences)"
// @Success 200 {object} ports.UserConfigResponse "Configuration successfully retrieved"
// @Failure 400 {object} gin.H "Invalid user ID or missing configType"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden: Cannot view another user's config"
// @Failure 404 {object} gin.H "ErrUserConfigNotFound"
// @Failure 500 {object} gin.H "Database error"
// @Router /users/{id}/config/{configType} [get]
func (h *UserConfigHandler) GetUserConfig(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	configType := c.Param(h.routes.ParamKeyConfigType)
	if configType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config type is required"})
		return
	}

	config, err := h.userConfigService.GetUserConfig(c.Request.Context(), targetUserID, configType)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapUserConfigToResponse(config))
}

// @Summary Delete User Configuration by Type
// @Description Deletes a specific configuration entry for a user. Enforces self-management.
// @Tags users, config
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param configType path string true "Configuration Type (e.g., notification_settings)"
// @Success 204 "Configuration successfully deleted (No Content)"
// @Failure 400 {object} gin.H "Invalid user ID or missing configType"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 403 {object} gin.H "Forbidden: Cannot delete another user's config"
// @Failure 404 {object} gin.H "ErrUserConfigNotFound"
// @Failure 500 {object} gin.H "Database error"
// @Router /users/{id}/config/{configType} [delete]
func (h *UserConfigHandler) DeleteUserConfig(c *gin.Context) {
	_, targetUserID, apiErr := h.checkUserOwnership(c)
	if apiErr != nil {
		c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
		return
	}

	configType := c.Param(h.routes.ParamKeyConfigType)
	if configType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config type is required"})
		return
	}

	err := h.userConfigService.DeleteUserConfig(c.Request.Context(), targetUserID, configType)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.Status(http.StatusNoContent)
}
