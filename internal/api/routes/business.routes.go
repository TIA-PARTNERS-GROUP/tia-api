package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupBusinessRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	
	businesses := api.Group(deps.Routes.BusinessBase)
	{
		
		businesses.GET("", deps.BusinessHandler.GetBusinesses)
		businesses.GET(deps.Routes.ParamID, deps.BusinessHandler.GetBusinessByID)
		
		businesses.GET(deps.Routes.BusinessConnects, deps.AuthMiddleware, deps.BusinessConnectionHandler.GetBusinessConnections)
		
		protectedBusinesses := businesses.Group("")
		protectedBusinesses.Use(deps.AuthMiddleware)
		{
			protectedBusinesses.POST("", deps.BusinessHandler.CreateBusiness)
			protectedBusinesses.PUT(deps.Routes.ParamID, deps.BusinessHandler.UpdateBusiness)
			protectedBusinesses.DELETE(deps.Routes.ParamID, deps.BusinessHandler.DeleteBusiness)
		}
		
		
		businessTags := businesses.Group(deps.Routes.BusinessTags)
		{
			businessTags.GET("", deps.BusinessTagHandler.GetBusinessTags)
			businessTags.POST("", deps.AuthMiddleware, deps.BusinessTagHandler.CreateBusinessTag)
		}
	}
}
