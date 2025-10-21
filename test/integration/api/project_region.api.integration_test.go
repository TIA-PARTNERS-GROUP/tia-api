package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_REGION_ID   = "AUS"
	TEST_REGION_NAME = "Australia"
)

// Helper to create a Region in the DB
func CreateTestRegion(t *testing.T) {
	region := models.Region{
		ID:   TEST_REGION_ID,
		Name: TEST_REGION_NAME,
	}
	result := testutil.TestDB.FirstOrCreate(&region)
	assert.NoError(t, result.Error, "Failed to create test region")
}

func TestProjectRegionAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	CreateTestRegion(t)

	// 1. Create Users
	managerUser, managerToken := CreateTestUserAndLogin(t, router, "region.manager@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "region.other@test.com", "ValidPass123!")

	// 2. Create Project
	project := CreateTestProjectHelper(t, router, managerUser, managerToken)
	projectID := project.ID

	regionsBaseURL := fmt.Sprintf("%s/projects/%d/regions", constants.AppRoutes.APIPrefix, projectID)
	regionSpecificURL := fmt.Sprintf("%s/%s", regionsBaseURL, TEST_REGION_ID)

	t.Run("Add Region - Forbidden (Not Manager)", func(t *testing.T) {
		addDTO := ports.AddProjectRegionInput{
			RegionID: TEST_REGION_ID,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, regionsBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Add Region - Success (Manager)", func(t *testing.T) {
		addDTO := ports.AddProjectRegionInput{
			ProjectID: 999, // Should be ignored
			RegionID:  TEST_REGION_ID,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, regionsBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var newRegion ports.ProjectRegionResponse
		json.Unmarshal(w.Body.Bytes(), &newRegion)
		assert.Equal(t, TEST_REGION_ID, newRegion.RegionID)
		assert.Equal(t, projectID, newRegion.ProjectID)
		assert.Equal(t, TEST_REGION_NAME, newRegion.Region.Name)
	})

	t.Run("Get Project Regions", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, regionsBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) // Any auth user
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []ports.ProjectRegionResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, len(resp))
		assert.Equal(t, TEST_REGION_ID, resp[0].RegionID)
	})

	t.Run("Remove Region - Forbidden (Not Manager)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, regionSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Remove Region - Success (Manager)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, regionSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Removal", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, regionsBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp []ports.ProjectRegionResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, len(resp))
	})
}
