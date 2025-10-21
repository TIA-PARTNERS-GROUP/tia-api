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
