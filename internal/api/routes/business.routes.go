package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupBusinessRoutes now takes *RouterDependencies
func SetupBusinessRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	// Use constants
	businesses := api.Group(deps.Routes.BusinessBase)
	{
		// Public GET routes
		businesses.GET("", deps.BusinessHandler.GetBusinesses)
		businesses.GET(deps.Routes.ParamID, deps.BusinessHandler.GetBusinessByID)

		// Specific middleware for getting connections
		businesses.GET(deps.Routes.BusinessConnects, deps.AuthMiddleware, deps.BusinessConnectionHandler.GetBusinessConnections)

		// Protected business actions
		protectedBusinesses := businesses.Group("")
		protectedBusinesses.Use(deps.AuthMiddleware)
		{
			protectedBusinesses.POST("", deps.BusinessHandler.CreateBusiness)
			protectedBusinesses.PUT(deps.Routes.ParamID, deps.BusinessHandler.UpdateBusiness)
			protectedBusinesses.DELETE(deps.Routes.ParamID, deps.BusinessHandler.DeleteBusiness)
		}

		// Business Tags sub-resource
		// This is now clearly the sub-route for tags
		businessTags := businesses.Group(deps.Routes.BusinessTags)
		{
			businessTags.GET("", deps.BusinessTagHandler.GetBusinessTags)
			businessTags.POST("", deps.AuthMiddleware, deps.BusinessTagHandler.CreateBusinessTag)
		}
	}
}
