package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupL2ERoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	l2e := api.Group(deps.Routes.L2EBase)
	l2e.Use(deps.AuthMiddleware)
	{
		l2e.POST("", deps.L2EHandler.CreateL2EResponse)
	}
}
