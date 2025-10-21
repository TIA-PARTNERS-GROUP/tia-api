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

		publications.GET(deps.Routes.PublicationByID, deps.PublicationHandler.GetPublicationByID)

		publications.GET(deps.Routes.PublicationBySlug, deps.PublicationHandler.GetPublicationBySlug)

		publications.PUT(deps.Routes.ParamID, deps.PublicationHandler.UpdatePublication)
		publications.DELETE(deps.Routes.ParamID, deps.PublicationHandler.DeletePublication)
	}
}
