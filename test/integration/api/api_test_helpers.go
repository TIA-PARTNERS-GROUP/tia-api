package main

import (
	"context"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	routes "github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/routes"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	userService := services.NewUserService(testutil.TestDB)
	authService := services.NewAuthService(testutil.TestDB)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			user, session, err := authService.ValidateToken(context.Background(), token)
			if err == nil {
				c.Set("user", user)
				c.Set("userID", user.ID)
				c.Set("sessionID", session.ID)
			}
		}
		c.Next()
	})

	routes.SetupRoutes(router, userHandler, authHandler)

	return router
}
