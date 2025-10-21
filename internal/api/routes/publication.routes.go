package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupPublicationRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	publications := api.Group(deps.Routes.PublicationBase)
	publications.Use(deps.AuthMiddleware)
	{
		publications.POST("", deps.PublicationHandler.CreatePublication)
		publications.GET("", deps.PublicationHandler.GetAllPublications)

		// Get by ID
		publications.GET(deps.Routes.PublicationByID, deps.PublicationHandler.GetPublicationByID)

		// Get by Slug (optional route if required, though ID is standard)
		publications.GET(deps.Routes.PublicationBySlug, deps.PublicationHandler.GetPublicationBySlug)

		// CRUD on a specific publication (requires ID in URL)
		publications.PUT(deps.Routes.ParamID, deps.PublicationHandler.UpdatePublication)
		publications.DELETE(deps.Routes.ParamID, deps.PublicationHandler.DeletePublication)
	}
}
