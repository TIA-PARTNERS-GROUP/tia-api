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
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT CONSTANTS
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// --- This part remains the same ---
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

	userHandler := handlers.NewUserHandler(userService)
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

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	authMiddlewareForTest := func(c *gin.Context) {
		// ... (your existing auth middleware is fine) ...
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

	// --- THIS IS THE MAIN CHANGE ---
	// 1. Create the RouterDependencies struct
	deps := &routes.RouterDependencies{
		AuthMiddleware:                authMiddlewareForTest,
		UserHandler:                   userHandler,
		AuthHandler:                   authHandler,
		BusinessHandler:               businessHandler,
		BusinessConnectionHandler:     businessConnectionHandler,
		BusinessTagHandler:            businessTagHandler,
		DailyActivityHandler:          dailyActivityHandler,
		DailyActivityEnrolmentHandler: dailyActivityEnrolmentHandler,
		EventHandler:                  eventHandler,
		FeedbackHandler:               feedbackHandler,
		InferredConnectionHandler:     inferredConnectionHandler,
		L2EHandler:                    l2eHandler,
		NotificationHandler:           notificationHandler,
		Routes:                        constants.AppRoutes, // <-- 2. Pass in the constants
	}

	// 3. Call the new, clean RegisterRoutes function
	routes.RegisterRoutes(router, deps)
	return router
}

func createJSONBody(t *testing.T, data interface{}) *bytes.Buffer {
	// ... (this helper is good, no change needed) ...
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return bytes.NewBuffer(jsonData)
}

func CreateTestUserAndLogin(t *testing.T, router *gin.Engine, email, password string) (models.User, string) {
	// ... (user creation is fine) ...
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

	// --- USE CONSTANTS FOR THE ROUTE ---
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
