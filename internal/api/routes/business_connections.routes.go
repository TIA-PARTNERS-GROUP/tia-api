package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupConnectionRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	connections := api.Group(deps.Routes.ConnectBase)
	connections.Use(deps.AuthMiddleware)
	{
		connections.POST("", deps.BusinessConnectionHandler.CreateBusinessConnection)
		connections.GET(deps.Routes.ParamID, deps.BusinessConnectionHandler.GetBusinessConnection)
		connections.PUT(deps.Routes.ParamID, deps.BusinessConnectionHandler.UpdateBusinessConnection)
		connections.POST(deps.Routes.ConnectAccept, deps.BusinessConnectionHandler.AcceptBusinessConnection)
		connections.POST(deps.Routes.ConnectReject, deps.BusinessConnectionHandler.RejectBusinessConnection)
		connections.DELETE(deps.Routes.ParamID, deps.BusinessConnectionHandler.DeleteBusinessConnection)
	}
}
