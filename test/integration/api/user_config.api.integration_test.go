package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest" // Corrected import (already in file, but used explicitly below)
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestUserConfigAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// 1. Create Users
	user, userToken := CreateTestUserAndLogin(t, router, "config.user@test.com", "ValidPass123!")
	// FIX 1: Removed otherToken as it was not used.
	otherUser, _ := CreateTestUserAndLogin(t, router, "config.other@test.com", "ValidPassword123!")

	const configType = "user_preferences"
	const configType2 = "notification_settings"

	// Base URL for the authenticated user's configs
	userConfigBaseURL := fmt.Sprintf("%s/users/%d/config", constants.AppRoutes.APIPrefix, user.ID)
	// Base URL for another user's configs (to test forbidden access)
	otherUserConfigBaseURL := fmt.Sprintf("%s/users/%d/config", constants.AppRoutes.APIPrefix, otherUser.ID)

	t.Run("Set Config (UPSERT) - Success", func(t *testing.T) {
		configData := `{"theme": "dark", "notifications_on": true}`
		setDTO := ports.SetUserConfigInput{
			ConfigType: configType,
			Config:     datatypes.JSON([]byte(configData)),
		}
		body, _ := json.Marshal(setDTO)
		req, _ := http.NewRequest(http.MethodPut, userConfigBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Config setting (UPSERT) failed")
		var resp ports.UserConfigResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, user.ID, resp.UserID)
		assert.Equal(t, configType, resp.ConfigType)
		assert.JSONEq(t, configData, string(resp.Config))
	})

	t.Run("Set Config - Forbidden (Targeting Another User)", func(t *testing.T) {
		configData := `{"theme": "light"}`
		setDTO := ports.SetUserConfigInput{
			ConfigType: configType,
			Config:     datatypes.JSON([]byte(configData)),
		}
		body, _ := json.Marshal(setDTO)
		req, _ := http.NewRequest(http.MethodPut, otherUserConfigBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken) // Authenticated as 'user', but targeting 'otherUser'

		// FIX 2: Corrected typo from htttest.NewRecorder() to httptest.NewRecorder()
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Get Config - Success", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", userConfigBaseURL, configType)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.UserConfigResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, user.ID, resp.UserID)
		assert.Equal(t, configType, resp.ConfigType)
		assert.Contains(t, string(resp.Config), "dark")
	})

	t.Run("Get Config - NotFound", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", userConfigBaseURL, configType2)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Update Config (UPSERT) - Success", func(t *testing.T) {
		configData := `{"theme": "light", "notifications_on": false, "new_field": 1}`
		setDTO := ports.SetUserConfigInput{
			ConfigType: configType,
			Config:     datatypes.JSON([]byte(configData)),
		}
		body, _ := json.Marshal(setDTO)
		req, _ := http.NewRequest(http.MethodPut, userConfigBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Config updating (UPSERT) failed")
		var resp ports.UserConfigResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, user.ID, resp.UserID)
		assert.JSONEq(t, configData, string(resp.Config))
	})

	t.Run("Delete Config - Success", func(t *testing.T) {
		// Set a second config type first
		setDTO := ports.SetUserConfigInput{
			ConfigType: configType2,
			Config:     datatypes.JSON([]byte(`{"volume": 50}`)),
		}
		body, _ := json.Marshal(setDTO)
		req, _ := http.NewRequest(http.MethodPut, userConfigBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+userToken)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Now delete it
		url := fmt.Sprintf("%s/%s", userConfigBaseURL, configType2)
		req, _ = http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", userConfigBaseURL, configType2)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
