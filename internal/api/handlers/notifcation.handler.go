package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NotificationHandler struct {
	service  *services.NotificationService
	validate *validator.Validate
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service:  service,
		validate: validator.New(),
	}
}

// @Summary      Create Notification
// @Description  Creates a new notification. SenderUserID can be null for system notifications. (Protected)
// @Tags         Notifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        notification body ports.CreateNotificationInput true "Notification Details"
// @Success      201 {object} ports.NotificationResponse
// @Failure      400 {object} map[string]string "Validation error or invalid request body"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      404 {object} map[string]string "Receiver user not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /notifications [post]
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var input ports.CreateNotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Basic validation based on DTO tags
	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: Add logic here to verify SenderUserID if provided,
	// e.g., check if it matches the authenticated user or if the user has permission.
	// For now, we trust the input or allow nil for system messages.

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

// @Summary      Get User Notifications
// @Description  Retrieves notifications for the specified user, optionally filtering by read status. (Protected)
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Param        read query bool false "Filter by read status (true or false)"
// @Success      200 {array} ports.NotificationResponse
// @Failure      400 {object} map[string]string "Invalid User ID or read parameter"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden - Cannot access another user's notifications"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{id}/notifications [get]
func (h *NotificationHandler) GetNotificationsForUser(c *gin.Context) {
	userIDParam := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// --- Authorization Check ---
	authUserIDVal, _ := c.Get("userID")
	authUserID, _ := authUserIDVal.(uint)
	// Simple check: User can only get their own notifications.
	// TODO: Implement admin role check if needed.
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// --- End Authorization Check ---

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

// @Summary      Mark Notification As Read
// @Description  Marks a specific notification as read for the authenticated user. (Protected)
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        userID path int true "User ID"
// @Param        notificationID path int true "Notification ID"
// @Success      200 {object} ports.NotificationResponse
// @Failure      400 {object} map[string]string "Invalid User ID or Notification ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden - Cannot modify another user's notification"
// @Failure      404 {object} map[string]string "Notification not found for this user"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{userID}/notifications/{notificationID}/read [patch]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDParam := c.Param("userID")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notificationIDParam := c.Param("notificationID")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// --- Authorization Check ---
	authUserIDVal, _ := c.Get("userID")
	authUserID, _ := authUserIDVal.(uint)
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// --- End Authorization Check ---

	notification, err := h.service.MarkAsRead(c.Request.Context(), uint(notificationID), authUserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotificationNotFound) {
			// Service returns this if ID doesn't exist OR user doesn't match
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, ports.MapNotificationToResponse(notification))
}

// @Summary      Mark All Notifications As Read
// @Description  Marks all unread notifications as read for the authenticated user. (Protected)
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        userID path int true "User ID"
// @Success      200 {object} map[string]int64 "Returns the count of notifications marked as read"
// @Failure      400 {object} map[string]string "Invalid User ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden - Cannot modify another user's notifications"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{userID}/notifications/read-all [patch]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDParam := c.Param("userID")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// --- Authorization Check ---
	authUserIDVal, _ := c.Get("userID")
	authUserID, _ := authUserIDVal.(uint)
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// --- End Authorization Check ---

	count, err := h.service.MarkAllAsRead(c.Request.Context(), authUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"marked_as_read_count": count})
}

// @Summary      Delete Notification
// @Description  Deletes a specific notification for the authenticated user. (Protected)
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        userID path int true "User ID"
// @Param        notificationID path int true "Notification ID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid User ID or Notification ID"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Failure      403 {object} map[string]string "Forbidden - Cannot delete another user's notification"
// @Failure      404 {object} map[string]string "Notification not found for this user"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /users/{userID}/notifications/{notificationID} [delete]
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userIDParam := c.Param("userID")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	notificationIDParam := c.Param("notificationID")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// --- Authorization Check ---
	authUserIDVal, _ := c.Get("userID")
	authUserID, _ := authUserIDVal.(uint)
	if authUserID != uint(targetUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// --- End Authorization Check ---

	err = h.service.DeleteNotification(c.Request.Context(), uint(notificationID), authUserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotificationNotFound) {
			// Service returns this if ID doesn't exist OR user doesn't match
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	c.Status(http.StatusNoContent)
}
