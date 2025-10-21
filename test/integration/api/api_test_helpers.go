package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/handlers"
	routes "github.com/TIA-PARTNERS-GROUP/tia-api/internal/api/routes"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {

	userService := services.NewUserService(testutil.TestDB)
	authService := services.NewAuthService(testutil.TestDB)
	businessService := services.NewBusinessService(testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	businessTagService := services.NewBusinessTagService(testutil.TestDB)
	dailyActivityService := services.NewDailyActivityService(testutil.TestDB)
	dailyActivityEnrolmentService := services.NewDailyActivityEnrolmentService(testutil.TestDB)
	eventService := services.NewEventService(testutil.TestDB)
	feedbackService := services.NewFeedbackService(testutil.TestDB)
	inferredConnectionService := services.NewInferredConnectionService(testutil.TestDB)
	l2eResponseService := services.NewL2EResponseService(testutil.TestDB)
	notificationService := services.NewNotificationService(testutil.TestDB)
	projectService := services.NewProjectService(testutil.TestDB)
	projectApplicantService := services.NewProjectApplicantService(testutil.TestDB) 
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)       
	projectRegionService := services.NewProjectRegionService(testutil.TestDB)       
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)         
	publicationService := services.NewPublicationService(testutil.TestDB)           
	skillService := services.NewSkillService(testutil.TestDB)                       
	subscriptionService := services.NewSubscriptionService(testutil.TestDB)         
	userSubscriptionService := services.NewUserSubscriptionService(testutil.TestDB) 
	userConfigService := services.NewUserConfigService(testutil.TestDB)             
	userSkillService := services.NewUserSkillService(testutil.TestDB)               

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

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	authMiddlewareForTest := func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		token := parts[1]
		user, session, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		c.Set(constants.AppRoutes.ContextKeyUser, user)
		c.Set(constants.AppRoutes.ContextKeyUserID, user.ID)
		c.Set(constants.AppRoutes.ContextKeySessionID, session.ID)
		c.Next()
	}

	deps := &routes.RouterDependencies{
		AuthMiddleware:                authMiddlewareForTest,
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

	routes.RegisterRoutes(router, deps)
	return router
}
func createJSONBody(t *testing.T, data interface{}) *bytes.Buffer {

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return bytes.NewBuffer(jsonData)
}
func CreateTestUserAndLogin(t *testing.T, router *gin.Engine, email, password string) (models.User, string) {

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	user := models.User{
		FirstName:    "Test",
		LoginEmail:   email,
		PasswordHash: &hashedPassword,
		Active:       true,
	}
	result := testutil.TestDB.Create(&user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)
	loginDTO := ports.LoginInput{
		LoginEmail: email,
		Password:   password,
	}
	body, err := json.Marshal(loginDTO)
	assert.NoError(t, err)

	loginPath := constants.AppRoutes.APIPrefix + constants.AppRoutes.AuthBase + constants.AppRoutes.Login
	req, _ := http.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Login request should be successful")
	var loginResponse ports.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResponse.Token, "Token should not be empty after login")
	return user, loginResponse.Token
}

func BoolPtr(val bool) *bool {
	return &val
}

func StrPtr(val string) *string {
	return &val
}

func IntPtr(val int) *int {
	return &val
}
