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

// @Summary Create New Subscription Plan
// @Description Creates a new recurring subscription plan definition. Requires authentication (implies admin/privileged access).
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan body ports.CreateSubscriptionInput true "Subscription plan details (Name, Price, ValidDays/ValidMonths)"
// @Success 201 {object} ports.SubscriptionResponse "Subscription plan created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "ErrSubscriptionNameExists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /subscriptions [post]
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

// @Summary Get Subscription Plan by ID
// @Description Retrieves a specific subscription plan definition by its unique ID.
// @Tags subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription Plan ID"
// @Success 200 {object} ports.SubscriptionResponse "Subscription plan retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid subscription ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrSubscriptionNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param(h.routes.ParamKeyID)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

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

// @Summary Subscribe User to a Plan
// @Description Creates a new UserSubscription record for the authenticated user, starting their access to a plan. Enforces self-subscription.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subscription body ports.UserSubscribeInput true "Subscription details (UserID and SubscriptionID)"
// @Success 201 {object} ports.UserSubscriptionResponse "User subscribed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation failed"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden: Cannot subscribe for another user"
// @Failure 404 {object} map[string]interface{} "ErrSubscriptionNotFound or ErrUserNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /subscriptions/subscribe [post]
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
	c.JSON(http.StatusCreated, ports.MapUserSubscriptionToResponse(userSubscription))
}
