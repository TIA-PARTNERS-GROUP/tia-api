package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupFeedbackRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	feedback := api.Group(deps.Routes.FeedbackBase)
	{
		// Public POST create
		feedback.POST("", deps.FeedbackHandler.CreateFeedback)

		// Protected GET/DELETE
		protectedFeedback := feedback.Group("")
		protectedFeedback.Use(deps.AuthMiddleware)
		{
			protectedFeedback.GET("", deps.FeedbackHandler.GetAllFeedback)
			protectedFeedback.GET(deps.Routes.ParamID, deps.FeedbackHandler.GetFeedbackByID)
			protectedFeedback.DELETE(deps.Routes.ParamID, deps.FeedbackHandler.DeleteFeedback)
		}
	}
}
