package api

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	authMiddleware := func(c *gin.Context) {
		c.Next()
	}

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)

			protectedAuth := auth.Group("/").Use(authMiddleware)
			{
				protectedAuth.POST("/logout", authHandler.Logout)
				protectedAuth.GET("/me", authHandler.GetCurrentUser)
			}
		}

		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
