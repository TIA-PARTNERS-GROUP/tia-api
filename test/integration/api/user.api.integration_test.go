package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	// "github.com/TIA-PARTNERS-GROUP/tia-api/internal/models" // No longer needed
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	// "github.com/TIA-PARTNERS-GROUP/tia-api/pkg/utils" // No longer needed
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

// In test/integration/api/user.api.integration_test.go

func TestUserAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	constApiPrefix := constants.AppRoutes.APIPrefix
	constUserBase := constApiPrefix + constants.AppRoutes.UsersBase

	// --- Create User to be Managed ---
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
	userIDToManage := createdUser.ID // Get the ID of the newly created user

	// --- Log in as the CREATED user to test self-update/delete ---
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
	token := loginResponse.Token // Use the token of the user being managed

	t.Run("Get User", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) // Get self
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
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) // Update self
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Expect 200 now
		var updatedUser ports.UserResponse
		json.Unmarshal(w.Body.Bytes(), &updatedUser)
		assert.Equal(t, "ApiUpdated", updatedUser.FirstName)
	})

	t.Run("Delete User", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage) // Delete self
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code) // Expect 204 now
	})

	t.Run("Verify Deletion", func(t *testing.T) {

		_, viewerToken := CreateTestUserAndLogin(t, router, "viewer@test.com", "ValidPass123!")

		url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+viewerToken) // Use the viewer's token
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Now that authentication passes, the handler's logic runs,
		// the service returns ErrUserNotFound, and the handler returns 404.
		assert.Equal(t, http.StatusNotFound, w.Code) // Expect 404 now

		// Option 2 (Less ideal but works if you don't need role separation):
		// Re-login as the original 'admin' user if you created one at the start.
		// However, the current test setup logs in as the user being managed,
		// whose token becomes invalid after deletion (or should).

		// --- Original Incorrect Logic ---
		// url := fmt.Sprintf("%s/%d", constUserBase, userIDToManage)
		// req, _ := http.NewRequest(http.MethodGet, url, nil) // No token!
		// w := httptest.NewRecorder()
		// router.ServeHTTP(w, req)
		// assert.Equal(t, http.StatusNotFound, w.Code) // This fails because it gets 401/403 first
	})

	// Optional: Add separate test cases for admin actions if you implement admin roles
}
