package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupNotificationRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	notifications := api.Group(deps.Routes.NotifyBase)
	notifications.Use(deps.AuthMiddleware)
	{
		notifications.POST("", deps.NotificationHandler.CreateNotification)
	}
}
