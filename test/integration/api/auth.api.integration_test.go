package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI_Integration_LoginLogout(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

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
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse ports.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NotEmpty(t, loginResponse.Token, "Token should not be empty")

	token := loginResponse.Token

	// 3. Test GetCurrentUser Endpoint (with valid token)
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var currentUser ports.UserResponse
	json.Unmarshal(w.Body.Bytes(), &currentUser)
	assert.Equal(t, user.LoginEmail, currentUser.LoginEmail)

	// 4. Test Logout Endpoint
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// 5. Verify token is now invalid
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
