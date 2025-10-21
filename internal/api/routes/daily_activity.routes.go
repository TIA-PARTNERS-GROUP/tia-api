package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupDailyActivityRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	dailyActivities := api.Group(deps.Routes.DailyActBase)
	{
		
		dailyActivities.GET("", deps.DailyActivityHandler.GetAllDailyActivities)
		dailyActivities.GET(deps.Routes.ParamID, deps.DailyActivityHandler.GetDailyActivityByID)
		
		dailyActivities.POST("", deps.AuthMiddleware, deps.DailyActivityHandler.CreateDailyActivity)
		
		enrolments := dailyActivities.Group(deps.Routes.DailyActEnrol)
		{
			enrolments.GET("", deps.DailyActivityEnrolmentHandler.GetEnrolmentsForActivity)
			
			enrolments.POST("", deps.AuthMiddleware, deps.DailyActivityEnrolmentHandler.EnrolUser)
			enrolments.DELETE("", deps.AuthMiddleware, deps.DailyActivityEnrolmentHandler.WithdrawUser)
		}
	}
}
