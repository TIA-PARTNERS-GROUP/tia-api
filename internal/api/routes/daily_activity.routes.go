package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupDailyActivityRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	dailyActivities := api.Group(deps.Routes.DailyActBase)
	{
		// Public GET routes
		dailyActivities.GET("", deps.DailyActivityHandler.GetAllDailyActivities)
		dailyActivities.GET(deps.Routes.ParamID, deps.DailyActivityHandler.GetDailyActivityByID)

		// Protected POST create
		dailyActivities.POST("", deps.AuthMiddleware, deps.DailyActivityHandler.CreateDailyActivity)

		// Enrolments sub-resource
		enrolments := dailyActivities.Group(deps.Routes.DailyActEnrol)
		{
			enrolments.GET("", deps.DailyActivityEnrolmentHandler.GetEnrolmentsForActivity)

			// Protected POST/DELETE enrol
			enrolments.POST("", deps.AuthMiddleware, deps.DailyActivityEnrolmentHandler.EnrolUser)
			enrolments.DELETE("", deps.AuthMiddleware, deps.DailyActivityEnrolmentHandler.WithdrawUser)
		}
	}
}
