package api

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/middleware"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authService *services.AuthService) {
	authMiddleware := middleware.AuthMiddleware(authService)

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

			protectedUsers := users.Group("/").Use(authMiddleware)
			{
				protectedUsers.GET("/", userHandler.GetAllUsers)
				protectedUsers.GET("/:id", userHandler.GetUserByID)
				protectedUsers.PUT("/:id", userHandler.UpdateUser)
				protectedUsers.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}
}
