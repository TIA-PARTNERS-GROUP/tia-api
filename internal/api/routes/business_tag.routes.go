package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupBusinessTagRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	
	tags := api.Group(deps.Routes.TagsBase)
	tags.Use(deps.AuthMiddleware)
	{
		tags.DELETE(deps.Routes.ParamID, deps.BusinessTagHandler.DeleteBusinessTag)
	}
}
