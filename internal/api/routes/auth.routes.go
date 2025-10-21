package routes
import (
	"github.com/gin-gonic/gin"
)
func SetupAuthRoutes(api *gin.RouterGroup, deps *RouterDependencies) {
	auth := api.Group(deps.Routes.AuthBase)
	{
		auth.POST(deps.Routes.Login, deps.AuthHandler.Login)
		auth.POST(deps.Routes.Logout, deps.AuthMiddleware, deps.AuthHandler.Logout)
		auth.GET(deps.Routes.Me, deps.AuthMiddleware, deps.AuthHandler.GetCurrentUser)
	}
}
