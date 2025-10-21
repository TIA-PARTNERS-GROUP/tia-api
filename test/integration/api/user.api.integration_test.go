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

func TestUserAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// --- USE HELPER to get admin token for protected routes ---
	_, token := CreateTestUserAndLogin(t, router, "admin@api.com", "ValidPassword123!")

	// --- USE CONSTANTS ---
	constUserBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase

	var createdUser ports.UserResponse

	t.Run("Create User", func(t *testing.T) {
		createDTO := ports.UserCreationSchema{
			FirstName:  "ApiTest",
			LastName:   "User",
			LoginEmail: "api@test.com",
			Password:   "ValidPassword123!",
		}
		body, _ := json.Marshal(createDTO)

		// Use constant path (this route is public)
		req, _ := http.NewRequest(http.MethodPost, constUserBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdUser)
		assert.Equal(t, "ApiTest", createdUser.FirstName)
		assert.NotZero(t, createdUser.ID)
	})

	t.Run("Get User", func(t *testing.T) {
		// Use constant path
		url := fmt.Sprintf("%s/%d", constUserBase, createdUser.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedUser ports.UserResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedUser)
		assert.Equal(t, createdUser.ID, fetchedUser.ID)
	})

	t.Run("Update User", func(t *testing.T) {
		updateDTO := `{"first_name": "ApiUpdated"}`
		// Use constant path
		url := fmt.Sprintf("%s/%d", constUserBase, createdUser.ID)
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
		// Use constant path
		url := fmt.Sprintf("%s/%d", constUserBase, createdUser.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		// Use constant path
		url := fmt.Sprintf("%s/%d", constUserBase, createdUser.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
