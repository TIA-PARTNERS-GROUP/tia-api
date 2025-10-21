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
			protectedUsers.GET(deps.Routes.UserApplications, deps.ProjectApplicantHandler.GetApplicationsForUser)
			protectedUsers.GET(deps.Routes.ProjectMemberships, deps.ProjectMemberHandler.GetProjectsByUser)

			// --- FIX: USE EXPLICIT, NON-NESTED ROUTES FOR SUBSCRIPTIONS ---
			// This matches the pattern of other nested GETs (like UserApplications)

			// 1. GET /users/:id/subscriptions (List active subscriptions)
			protectedUsers.GET(deps.Routes.UserSubscriptions, deps.UserSubscriptionHandler.GetSubscriptionsForUser)

			// 2. DELETE /users/:id/subscriptions/:userSubscriptionID (Cancel one subscription)
			// Note: UserSubscriptionCancel must be defined as /:id/subscriptions/:userSubscriptionID
			protectedUsers.DELETE(deps.Routes.UserSubscriptionCancel, deps.UserSubscriptionHandler.CancelSubscription)

			// --- END FIX ---

			protectedUsers.PUT(deps.Routes.ParamID, deps.UserHandler.UpdateUser)
			protectedUsers.DELETE(deps.Routes.ParamID, deps.UserHandler.DeleteUser)

			userConfig := protectedUsers.Group(deps.Routes.UserConfigBase) // /users/:id/config
			{
				userConfig.PUT("", deps.UserConfigHandler.SetUserConfig)                                // PUT /users/:id/config
				userConfig.GET(deps.Routes.ParamConfigType, deps.UserConfigHandler.GetUserConfig)       // GET /users/:id/config/:configType
				userConfig.DELETE(deps.Routes.ParamConfigType, deps.UserConfigHandler.DeleteUserConfig) // DELETE /users/:id/config/:configType
			}

			userSkills := protectedUsers.Group(deps.Routes.UserSkillsBase)
			{
				userSkills.POST("", deps.UserSkillHandler.AddUserSkill)
				userSkills.GET("", deps.UserSkillHandler.GetUserSkills)
				userSkills.PUT(deps.Routes.ParamSkillID, deps.UserSkillHandler.UpdateUserSkill)
				userSkills.DELETE(deps.Routes.ParamSkillID, deps.UserSkillHandler.RemoveUserSkill)
			}

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
