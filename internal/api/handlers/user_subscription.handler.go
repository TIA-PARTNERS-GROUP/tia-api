package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
)

type UserSubscriptionHandler struct {
	userSubscriptionService *services.UserSubscriptionService
	routes                  *constants.Routes
}

func NewUserSubscriptionHandler(userSubscriptionService *services.UserSubscriptionService, routes *constants.Routes) *UserSubscriptionHandler {
	return &UserSubscriptionHandler{
		userSubscriptionService: userSubscriptionService,
		routes:                  routes,
	}
}

// getAuthUserID retrieves the authenticated user's ID from the context.
func (h *UserSubscriptionHandler) getAuthUserID(c *gin.Context) (uint, error) {
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

// GetSubscriptionsForUser handles GET /users/:id/subscriptions
func (h *UserSubscriptionHandler) GetSubscriptionsForUser(c *gin.Context) {
	targetUserIDStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Auth: Users can only view their own subscriptions
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only view your own subscriptions"})
		return
	}

	subscriptions, err := h.userSubscriptionService.GetSubscriptionsForUser(c.Request.Context(), authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	responses := make([]ports.UserSubscriptionResponse, len(subscriptions))
	for i, sub := range subscriptions {
		// NOTE: The service must ensure `Subscription` is preloaded for mapping to work
		responses[i] = ports.MapUserSubscriptionToResponse(&sub)
	}

	c.JSON(http.StatusOK, responses)
}

// CancelSubscription handles DELETE /users/:userID/subscriptions/:subID
func (h *UserSubscriptionHandler) CancelSubscription(c *gin.Context) {
	targetUserIDStr := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userSubIDStr := c.Param(h.routes.ParamKeySubscriptionID)
	userSubID, err := strconv.ParseUint(userSubIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user subscription ID"})
		return
	}

	authUserID, err := h.getAuthUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Auth: Users can only cancel their own subscription records
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only cancel your own subscriptions"})
		return
	}

	// Step 1: Verify ownership of the subscription record
	userSub, err := h.userSubscriptionService.GetUserSubscriptionByID(c.Request.Context(), uint(userSubID))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving subscription record"})
		return
	}

	if userSub.UserID != authUserID {
		// Found the record, but it belongs to someone else. Return Forbidden.
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Subscription record does not belong to you"})
		return
	}

	// Step 2: Delete/Cancel the record
	err = h.userSubscriptionService.CancelSubscription(c.Request.Context(), uint(userSubID))
	if err != nil {
		var apiErr *ports.ApiError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred during cancellation"})
		return
	}

	c.Status(http.StatusNoContent)
}
