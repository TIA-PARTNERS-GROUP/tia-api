package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes configures top-level notification routes (e.g., create)
// Note: GET/PATCH/DELETE routes are under /users/:id/notifications
func SetupNotificationRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	notifications := api.Group(deps.Routes.NotifyBase)
	notifications.Use(deps.AuthMiddleware)
	{
		notifications.POST("", deps.NotificationHandler.CreateNotification)
	}
}
