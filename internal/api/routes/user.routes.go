package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	users := api.Group(deps.Routes.UsersBase)
	{
		// Public routes
		users.POST("", deps.UserHandler.CreateUser)
		users.GET(deps.Routes.UserEnrolments, deps.DailyActivityEnrolmentHandler.GetEnrolmentsForUser)
		users.GET(deps.Routes.UserL2EResponses, deps.L2EHandler.GetL2EResponsesForUser)
		// users.GET("/:id/idea-votes", ideaVoteHandler.GetIdeaVotesByUser) // If needed

		// Protected routes group
		protectedUsers := users.Group("")
		protectedUsers.Use(deps.AuthMiddleware)
		{
			// GET /users
			protectedUsers.GET("", deps.UserHandler.GetAllUsers)
			// GET /users/:id
			protectedUsers.GET(deps.Routes.ParamID, deps.UserHandler.GetUserByID)
			// PUT /users/:id
			protectedUsers.PUT(deps.Routes.ParamID, deps.UserHandler.UpdateUser)
			// DELETE /users/:id
			protectedUsers.DELETE(deps.Routes.ParamID, deps.UserHandler.DeleteUser)

			// Notification Routes (nested)
			// Base: /users/:id/notifications
			notifications := protectedUsers.Group(deps.Routes.UserNotifications)
			{
				// GET /users/:id/notifications
				notifications.GET("", deps.NotificationHandler.GetNotificationsForUser)
				// PATCH /users/:id/notifications/read-all
				notifications.PATCH(deps.Routes.UserNotifyReadAll, deps.NotificationHandler.MarkAllAsRead)
				// PATCH /users/:id/notifications/:notificationID/read
				notifications.PATCH(deps.Routes.UserNotifyReadOne, deps.NotificationHandler.MarkAsRead)
				// DELETE /users/:id/notifications/:notificationID
				notifications.DELETE(deps.Routes.ParamNotificationID, deps.NotificationHandler.DeleteNotification)
			}
		}
	}
}
