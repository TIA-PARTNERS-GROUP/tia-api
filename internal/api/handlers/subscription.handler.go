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

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
	validate            *validator.Validate
	routes              *constants.Routes
}

func NewSubscriptionHandler(subscriptionService *services.SubscriptionService, routes *constants.Routes) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validate:            validator.New(),
		routes:              routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *SubscriptionHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

// CreateSubscription handles POST /subscriptions
// NOTE: Assuming this is an admin-level endpoint, but only checking for basic authentication
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.CreateSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.subscriptionService.CreateSubscription(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapSubscriptionToResponse(subscription))
}

// GetSubscriptionByID handles GET /subscriptions/:id
func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	// Auth: Any authenticated user can view a subscription plan
	if _, err := h.getAuthUserID(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	subscription, err := h.subscriptionService.GetSubscriptionByID(c.Request.Context(), uint(id))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, ports.MapSubscriptionToResponse(subscription))
}

// SubscribeUser handles POST /subscriptions/subscribe
func (h *SubscriptionHandler) SubscribeUser(c *gin.Context) {
	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input ports.UserSubscribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure the user ID in the body matches the authenticated user ID (self-subscribe)
	if input.UserID != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only subscribe for yourself"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userSubscription, err := h.subscriptionService.SubscribeUser(c.Request.Context(), input)
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	// NOTE: MapUserSubscriptionToResponse is assumed to exist in ports package
	c.JSON(http.StatusCreated, ports.MapUserSubscriptionToResponse(userSubscription))
}
