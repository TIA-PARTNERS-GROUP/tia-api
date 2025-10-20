package api

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUserByID)
		}
	}
}
