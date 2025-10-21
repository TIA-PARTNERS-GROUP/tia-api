package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupBusinessTagRoutes configures top-level tag routes
func SetupBusinessTagRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	// This is now clearly the *top-level* tag route
	tags := api.Group(deps.Routes.TagsBase)
	tags.Use(deps.AuthMiddleware)
	{
		tags.DELETE(deps.Routes.ParamID, deps.BusinessTagHandler.DeleteBusinessTag)
	}
}
