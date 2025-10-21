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

type NotificationHandler struct {
	service  *services.NotificationService
	validate *validator.Validate
	routes   *constants.Routes
}

func NewNotificationHandler(service *services.NotificationService, routes *constants.Routes) *NotificationHandler {
	return &NotificationHandler{
		service:  service,
		validate: validator.New(),
		routes:   routes,
	}
}

// @Summary Create Notification (Internal/System Use)
// @Description Creates a new notification record. Requires authentication, typically used by system services.
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body ports.CreateNotificationInput true "Notification details (ReceiverUserID, Title, Message, Type)"
// @Success 201 {object} ports.NotificationResponse "Notification created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "ErrReceiverNotFound"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {

	_, exists := c.Get(h.routes.ContextKeyUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var input ports.CreateNotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notification, err := h.service.CreateNotification(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ports.ErrReceiverNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}
	c.JSON(http.StatusCreated, ports.MapNotificationToResponse(notification))
}

// @Summary Get Notifications for User
// @Description Retrieves all notifications for the specified user. Requires authentication and self-management.
// @Tags notifications, users
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID"
// @Param read query bool false "Filter by read status (true/false)"
// @Success 200 {array} ports.NotificationResponse "List of notifications"
// @Failure 400 {object} map[string]interface{} "Invalid user ID or query parameter"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the target user)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/notifications [get]
func (h *NotificationHandler) GetNotificationsForUser(c *gin.Context) {

	userIDParam := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var readFilter *bool
	readQuery := c.Query("read")
	if readQuery != "" {
		readVal, err := strconv.ParseBool(readQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'read' query parameter"})
			return
		}
		readFilter = &readVal
	}
	notifications, err := h.service.GetNotificationsForUser(c.Request.Context(), uint(targetUserID), readFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}
	response := make([]ports.NotificationResponse, len(notifications))
	for i, n := range notifications {
		response[i] = ports.MapNotificationToResponse(&n)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary Mark Single Notification as Read
// @Description Marks a specific notification as read. Requires authentication and self-management.
// @Tags notifications, users
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID"
// @Param notificationID path int true "Notification ID"
// @Success 200 {object} ports.NotificationResponse "Notification marked as read"
// @Failure 400 {object} map[string]interface{} "Invalid user/notification ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the target user)"
// @Failure 404 {object} map[string]interface{} "Notification not found for this user"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/notifications/{notificationID}/read [patch]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {

	userIDParam := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notificationIDParam := c.Param(h.routes.ParamKeyNotificationID)
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	notification, err := h.service.MarkAsRead(c.Request.Context(), uint(notificationID), authUserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotificationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}
	c.JSON(http.StatusOK, ports.MapNotificationToResponse(notification))
}

// @Summary Mark All Notifications as Read
// @Description Marks all unread notifications for the specified user as read. Requires authentication and self-management.
// @Tags notifications, users
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID"
// @Success 200 {object} map[string]interface{} "Count of notifications marked as read"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the target user)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/notifications/read-all [patch]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {

	userIDParam := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	count, err := h.service.MarkAllAsRead(c.Request.Context(), authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"marked_as_read_count": count})
}

// @Summary Delete Single Notification
// @Description Deletes a specific notification for the user. Requires authentication and self-management.
// @Tags notifications, users
// @Produce json
// @Security BearerAuth
// @Param id path int true "Target User ID"
// @Param notificationID path int true "Notification ID"
// @Success 204 "Notification deleted successfully (No Content)"
// @Failure 400 {object} map[string]interface{} "Invalid user/notification ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden (Not the target user)"
// @Failure 404 {object} map[string]interface{} "Notification not found for this user"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/notifications/{notificationID} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {

	userIDParam := c.Param(h.routes.ParamKeyID)
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notificationIDParam := c.Param(h.routes.ParamKeyNotificationID)
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	authUserIDVal, _ := c.Get(h.routes.ContextKeyUserID)
	authUserID, ok := authUserIDVal.(uint)
	if !ok || authUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication context"})
		return
	}
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = h.service.DeleteNotification(c.Request.Context(), uint(notificationID), authUserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotificationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}
	c.Status(http.StatusNoContent)
}
