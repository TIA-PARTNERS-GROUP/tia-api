package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupSubscriptionRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	subs := api.Group(deps.Routes.SubscriptionBase)
	subs.Use(deps.AuthMiddleware)
	{
		subs.POST("", deps.SubscriptionHandler.CreateSubscription)
		subs.GET(deps.Routes.ParamID, deps.SubscriptionHandler.GetSubscriptionByID)

		subs.POST(deps.Routes.SubscriptionSubscribe, deps.SubscriptionHandler.SubscribeUser)
	}
}
