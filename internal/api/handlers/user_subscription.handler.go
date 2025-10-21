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

// @Summary Get Active Subscriptions for User
// @Description Retrieves all currently active user subscription records (those not yet expired). Enforces self-management.
// @Tags users, subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID (must match authenticated user)"
// @Success 200 {array} ports.UserSubscriptionResponse "List of active user subscriptions"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Cannot view another user's subscriptions"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/subscriptions [get]
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
		responses[i] = ports.MapUserSubscriptionToResponse(&sub)
	}

	c.JSON(http.StatusOK, responses)
}

// @Summary Cancel User Subscription
// @Description Cancels a specific user subscription record by deleting it from the database. Enforces self-management and ownership of the record.
// @Tags users, subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID (must match authenticated user)"
// @Param userSubscriptionID path int true "User Subscription Record ID to cancel"
// @Success 204 "Subscription cancelled successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid user or subscription ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: You are not the owner of this record"
// @Failure 404 {object} map[string]interface{} "ErrUserSubscriptionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal error during cancellation"
// @Router /users/{id}/subscriptions/{userSubscriptionID} [delete]
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

	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You can only cancel your own subscriptions"})
		return
	}

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
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Subscription record does not belong to you"})
		return
	}

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
