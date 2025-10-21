package main

import (
	"log"

	"github.com/TIA-PARTNERS-GROUP/tia-api/configs"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	handler "github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/middleware"
	routes "github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/routes"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/TIA-PARTNERS-GROUP/tia-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           TIA API
// @version         1.0
// @description     This is the API server for the TIA platform.
// @license.name  Apache 2.0
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @bearerFormat JWT
// @alias datatypes.JSON = interface{}
func main() {
	config := configs.LoadConfig()

	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.Business{},
		&models.Project{},
		&models.Skill{},
		&models.Publication{},
		&models.Notification{},
		&models.UserSkill{},
		&models.ProjectSkill{},
		&models.ProjectMember{},
		&models.BusinessConnection{},
		&models.BusinessTag{},
		&models.Feedback{},
		&models.ProjectApplicant{},
		&models.DailyActivity{},
		&models.DailyActivityEnrolment{},
		&models.UserDailyActivityProgress{},
		&models.Event{},
		&models.Subscription{},
		&models.UserSubscription{},
		&models.UserConfig{},
		&models.L2EResponse{},
		&models.Region{},
		&models.ProjectRegion{},
		&models.InferredConnection{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration successful.")

	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	businessService := services.NewBusinessService(db)
	businessConnectionService := services.NewBusinessConnectionService(db)
	businessTagService := services.NewBusinessTagService(db)
	dailyActivityService := services.NewDailyActivityService(db)
	dailyActivityEnrolmentService := services.NewDailyActivityEnrolmentService(db)
	eventService := services.NewEventService(db)
	feedbackService := services.NewFeedbackService(db)
	inferredConnectionService := services.NewInferredConnectionService(db)
	l2eResponseService := services.NewL2EResponseService(db)
	notificationService := services.NewNotificationService(db)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, &constants.AppRoutes)
	businessHandler := handlers.NewBusinessHandler(businessService)
	businessConnectionHandler := handlers.NewBusinessConnectionHandler(businessConnectionService)
	businessTagHandler := handlers.NewBusinessTagHandler(businessTagService)
	dailyActivityHandler := handlers.NewDailyActivityHandler(dailyActivityService)
	dailyActivityEnrolmentHandler := handlers.NewDailyActivityEnrolmentHandler(dailyActivityEnrolmentService)
	eventHandler := handlers.NewEventHandler(eventService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)
	inferredConnectionHandler := handlers.NewInferredConnectionHandler(inferredConnectionService)
	l2eHandler := handlers.NewL2EHandler(l2eResponseService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	authMiddleware := middleware.AuthMiddleware(authService)

	router := gin.Default()
	routes.RegisterRoutes(
		router,
		authMiddleware,
		userHandler,
		authHandler,
		businessHandler,
		businessConnectionHandler,
		businessTagHandler,
		dailyActivityHandler,
		dailyActivityEnrolmentHandler,
		eventHandler,
		feedbackHandler,
		inferredConnectionHandler,
		l2eHandler,
		notificationHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Starting server on http://localhost:8080")
	log.Println("Swagger UI available on http://localhost:8080/swagger/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
