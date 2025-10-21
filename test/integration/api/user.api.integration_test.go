package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" 
	
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestUserAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	constApiPrefix := constants.AppRoutes.APIPrefix
	constUserBase := constApiPrefix + constants.AppRoutes.UsersBase
	
	createDTO := ports.UserCreationSchema{
		FirstName:  "ApiTest",
		LastName:   "User",
		LoginEmail: "api@test.com",
		Password:   "ValidPassword123!",
	}
	createBody, _ := json.Marshal(createDTO)
	reqCreate, _ := http.NewRequest(http.MethodPost, constUserBase, bytes.NewBuffer(createBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, wCreate.Code, "User creation failed")
	var createdUser ports.UserResponse
	json.Unmarshal(wCreate.Body.Bytes(), &createdUser)
	assert.NotZero(t, createdUser.ID)
	userIDToManage := createdUser.ID 
	
	loginDTO := ports.LoginInput{LoginEmail: createDTO.LoginEmail, Password: createDTO.Password}
	loginBody, _ := json.Marshal(loginDTO)
	loginPath := constApiPrefix + constants.AppRoutes.AuthBase + constants.AppRoutes.Login
	reqLogin, _ := http.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	router.ServeHTTP(wLogin, reqLogin)
	assert.Equal(t, http.StatusOK, wLogin.Code, "Login failed for created user")
	var loginResponse ports.LoginResponse
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	token := loginResponse.Token 
	t.Run("Get User", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) 
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedUser ports.UserResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedUser)
		assert.Equal(t, userIDToManage, fetchedUser.ID)
	})
	t.Run("Update User", func(t *testing.T) {
		updateDTO := `{"first_name": "ApiUpdated"}`
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) 
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) 
		var updatedUser ports.UserResponse
		json.Unmarshal(w.Body.Bytes(), &updatedUser)
		assert.Equal(t, "ApiUpdated", updatedUser.FirstName)
	})
	t.Run("Delete User", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) 
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code) 
	})
	t.Run("Verify Deletion", func(t *testing.T) {
		_, viewerToken := CreateTestUserAndLogin(t, router, "viewer@test.com", "ValidPass123!")
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+viewerToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		
		assert.Equal(t, http.StatusNotFound, w.Code) 
		
		
		
		
		
		
		
		
		
		
	})
	
}
