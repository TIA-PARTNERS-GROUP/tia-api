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
