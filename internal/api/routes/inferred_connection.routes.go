package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupInferredConnectionRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	inferred := api.Group(deps.Routes.InferredBase)
	inferred.Use(deps.AuthMiddleware)
	{
		inferred.POST("", deps.InferredConnectionHandler.CreateInferredConnection)
		inferred.GET(deps.Routes.InferredBySource, deps.InferredConnectionHandler.GetConnectionsForSource)
	}
}
