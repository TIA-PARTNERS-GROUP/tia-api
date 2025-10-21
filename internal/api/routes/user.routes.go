package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupUserRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	users := api.Group(deps.Routes.UsersBase)
	{
		
		users.POST("", deps.UserHandler.CreateUser)
		users.GET(deps.Routes.UserEnrolments, deps.DailyActivityEnrolmentHandler.GetEnrolmentsForUser)
		users.GET(deps.Routes.UserL2EResponses, deps.L2EHandler.GetL2EResponsesForUser)
		
		
		protectedUsers := users.Group("")
		protectedUsers.Use(deps.AuthMiddleware)
		{
			
			protectedUsers.GET("", deps.UserHandler.GetAllUsers)
			
			protectedUsers.GET(deps.Routes.ParamID, deps.UserHandler.GetUserByID)
			
			protectedUsers.PUT(deps.Routes.ParamID, deps.UserHandler.UpdateUser)
			
			protectedUsers.DELETE(deps.Routes.ParamID, deps.UserHandler.DeleteUser)
			
			
			notifications := protectedUsers.Group(deps.Routes.UserNotifications)
			{
				
				notifications.GET("", deps.NotificationHandler.GetNotificationsForUser)
				
				notifications.PATCH(deps.Routes.UserNotifyReadAll, deps.NotificationHandler.MarkAllAsRead)
				
				notifications.PATCH(deps.Routes.UserNotifyReadOne, deps.NotificationHandler.MarkAsRead)
				
				notifications.DELETE(deps.Routes.ParamNotificationID, deps.NotificationHandler.DeleteNotification)
			}
		}
	}
}
