package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupConnectionRecommendationRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	connections := api.Group("/connections")
	{
		connections.GET("/complementary/:"+deps.Routes.ParamKeyUserID, deps.AuthMiddleware, deps.ConnectionHandler.GetComplementaryPartners)
		connections.GET("/alliance/:"+deps.Routes.ParamKeyUserID, deps.AuthMiddleware, deps.ConnectionHandler.GetAlliancePartners)
		connections.GET("/mastermind/:"+deps.Routes.ParamKeyUserID, deps.AuthMiddleware, deps.ConnectionHandler.GetMastermindPartners)
		connections.GET("/recommendations/:"+deps.Routes.ParamKeyUserID, deps.AuthMiddleware, deps.ConnectionHandler.GetAllRecommendations)
		connections.GET("/analysis/:"+deps.Routes.ParamKeyUserID, deps.AuthMiddleware, deps.ConnectionHandler.GetConnectionAnalysis)
	}
}
