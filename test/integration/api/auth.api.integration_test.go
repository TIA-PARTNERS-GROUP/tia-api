package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI_Integration_LoginLogout(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// --- USE CONSTANTS ---
	constApiPrefix := constants.AppRoutes.APIPrefix
	constAuthBase := constApiPrefix + constants.AppRoutes.AuthBase
	constLoginPath := constAuthBase + constants.AppRoutes.Login
	constMePath := constAuthBase + constants.AppRoutes.Me
	constLogoutPath := constAuthBase + constants.AppRoutes.Logout

	// 1. Seed a test user
	password := "ValidPassword123!"
	hashedPassword, _ := utils.HashPassword(password)
	user := models.User{FirstName: "ApiAuth", LoginEmail: "apiauth@test.com", PasswordHash: &hashedPassword, Active: true}
	testutil.TestDB.Select("*").Create(&user)

	// 2. Test Login Endpoint
	loginDTO := ports.LoginInput{
		LoginEmail: "apiauth@test.com",
		Password:   password,
	}
	body, _ := json.Marshal(loginDTO)
	// Use constant path
	req, _ := http.NewRequest(http.MethodPost, constLoginPath, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse ports.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NotEmpty(t, loginResponse.Token, "Token should not be empty")

	token := loginResponse.Token

	// 3. Test GetCurrentUser Endpoint (with valid token)
	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constMePath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var currentUser ports.UserResponse
	json.Unmarshal(w.Body.Bytes(), &currentUser)
	assert.Equal(t, user.LoginEmail, currentUser.LoginEmail)

	// 4. Test Logout Endpoint
	// Use constant path
	req, _ = http.NewRequest(http.MethodPost, constLogoutPath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// 5. Verify token is now invalid
	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constMePath, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
