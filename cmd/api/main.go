package main

import (
	"log"

	"github.com/TIA-PARTNERS-GROUP/tia-api/configs"
	_ "github.com/TIA-PARTNERS-GROUP/tia-api/docs"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/middleware"
	routes "github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/routes"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	projectService := services.NewProjectService(db)
	projectApplicantService := services.NewProjectApplicantService(db)
	projectMemberService := services.NewProjectMemberService(db)       
	projectRegionService := services.NewProjectRegionService(db)       
	projectSkillService := services.NewProjectSkillService(db)         
	publicationService := services.NewPublicationService(db)           
	skillService := services.NewSkillService(db)                       
	subscriptionService := services.NewSubscriptionService(db)         
	userSubscriptionService := services.NewUserSubscriptionService(db) 
	userConfigService := services.NewUserConfigService(db)             
	userSkillService := services.NewUserSkillService(db)               

	userHandler := handlers.NewUserHandler(userService, &constants.AppRoutes)
	authHandler := handlers.NewAuthHandler(authService, &constants.AppRoutes)
	businessHandler := handlers.NewBusinessHandler(businessService, &constants.AppRoutes)
	businessConnectionHandler := handlers.NewBusinessConnectionHandler(businessConnectionService, &constants.AppRoutes)
	businessTagHandler := handlers.NewBusinessTagHandler(businessTagService, &constants.AppRoutes)
	dailyActivityHandler := handlers.NewDailyActivityHandler(dailyActivityService, &constants.AppRoutes)
	dailyActivityEnrolmentHandler := handlers.NewDailyActivityEnrolmentHandler(dailyActivityEnrolmentService, &constants.AppRoutes)
	eventHandler := handlers.NewEventHandler(eventService, &constants.AppRoutes)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService, &constants.AppRoutes)
	inferredConnectionHandler := handlers.NewInferredConnectionHandler(inferredConnectionService, &constants.AppRoutes)
	l2eHandler := handlers.NewL2EHandler(l2eResponseService, &constants.AppRoutes)
	notificationHandler := handlers.NewNotificationHandler(notificationService, &constants.AppRoutes)
	projectHandler := handlers.NewProjectHandler(projectService, &constants.AppRoutes)
	projectApplicantHandler := handlers.NewProjectApplicantHandler(projectApplicantService, projectService, &constants.AppRoutes) 
	projectMemberHandler := handlers.NewProjectMemberHandler(projectMemberService, projectService, &constants.AppRoutes)          
	projectRegionHandler := handlers.NewProjectRegionHandler(projectRegionService, projectService, &constants.AppRoutes)          
	projectSkillHandler := handlers.NewProjectSkillHandler(projectSkillService, projectService, &constants.AppRoutes)             
	publicationHandler := handlers.NewPublicationHandler(publicationService, &constants.AppRoutes)                                
	skillHandler := handlers.NewSkillHandler(skillService, &constants.AppRoutes)                                                  
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, &constants.AppRoutes)                             
	userSubscriptionHandler := handlers.NewUserSubscriptionHandler(userSubscriptionService, &constants.AppRoutes)                 
	userConfigHandler := handlers.NewUserConfigHandler(userConfigService, &constants.AppRoutes)                                   
	userSkillHandler := handlers.NewUserSkillHandler(userSkillService, &constants.AppRoutes)                                      

	authMiddleware := middleware.AuthMiddleware(authService, &constants.AppRoutes)

	deps := &routes.RouterDependencies{
		AuthMiddleware:                authMiddleware,
		UserHandler:                   userHandler,
		AuthHandler:                   authHandler,
		BusinessHandler:               businessHandler,
		ProjectHandler:                projectHandler,
		BusinessConnectionHandler:     businessConnectionHandler,
		BusinessTagHandler:            businessTagHandler,
		DailyActivityHandler:          dailyActivityHandler,
		DailyActivityEnrolmentHandler: dailyActivityEnrolmentHandler,
		EventHandler:                  eventHandler,
		FeedbackHandler:               feedbackHandler,
		InferredConnectionHandler:     inferredConnectionHandler,
		L2EHandler:                    l2eHandler,
		NotificationHandler:           notificationHandler,
		ProjectApplicantHandler:       projectApplicantHandler, 
		ProjectMemberHandler:          projectMemberHandler,    
		ProjectRegionHandler:          projectRegionHandler,    
		ProjectSkillHandler:           projectSkillHandler,     
		PublicationHandler:            publicationHandler,      
		SkillHandler:                  skillHandler,            
		SubscriptionHandler:           subscriptionHandler,     
		UserSubscriptionHandler:       userSubscriptionHandler, 
		UserConfigHandler:             userConfigHandler,       
		UserSkillHandler:              userSkillHandler,        
		Routes:                        constants.AppRoutes,
	}
	router := gin.Default()

	routes.RegisterRoutes(router, deps)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("Starting server on http:
	log.Println("Swagger UI available on http:
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
