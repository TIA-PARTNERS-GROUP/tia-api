package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	events := api.Group(deps.Routes.EventBase)
	{
		// Public GET routes
		events.GET("", deps.EventHandler.GetEvents)
		events.GET(deps.Routes.ParamID, deps.EventHandler.GetEventByID)

		// Protected POST create
		events.POST("", deps.AuthMiddleware, deps.EventHandler.CreateEvent)
	}
}
