package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" // <-- IMPORT
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestInferredConnectionAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// 1. Setup: Create entities (e.g., user and business) and get token
	user, token := CreateTestUserAndLogin(t, router, "infer-user@test.com", "ValidPass123!")

	// Create a business to connect to
	biz := models.Business{Name: "Infer Target Biz", OperatorUserID: user.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Create(&biz)
	assert.NotZero(t, biz.ID)

	// --- USE CONSTANTS ---
	constInferredBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.InferredBase

	var createdConnection ports.InferredConnectionResponse

	t.Run("Create Inferred Connection", func(t *testing.T) {
		createDTO := ports.CreateInferredConnectionInput{
			SourceEntityType: "user",
			SourceEntityID:   user.ID,
			TargetEntityType: "business",
			TargetEntityID:   biz.ID,
			ConnectionType:   "Potential_Partner",
			ConfidenceScore:  0.85,
			ModelVersion:     "v1.0",
		}
		body, _ := json.Marshal(createDTO)
		// Use constant path
		req, _ := http.NewRequest(http.MethodPost, constInferredBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdConnection)
		assert.NotZero(t, createdConnection.ID)
		assert.Equal(t, "user", createdConnection.SourceEntityType)
		assert.Equal(t, user.ID, createdConnection.SourceEntityID)
		assert.Equal(t, "business", createdConnection.TargetEntityType)
		assert.Equal(t, biz.ID, createdConnection.TargetEntityID)
		assert.Equal(t, "Potential_Partner", createdConnection.ConnectionType)
		assert.Equal(t, 0.85, createdConnection.ConfidenceScore)
	})

	t.Run("Get Connections For Source", func(t *testing.T) {
		// Build path safely from constants
		getSourcePath := strings.Replace(constants.AppRoutes.InferredBySource, ":entityType", "user", 1)
		getSourcePath = strings.Replace(getSourcePath, ":entityID", fmt.Sprintf("%d", user.ID), 1)
		url := constInferredBase + getSourcePath

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var connections []ports.InferredConnectionResponse
		json.Unmarshal(w.Body.Bytes(), &connections)

		assert.Len(t, connections, 1)
		if len(connections) > 0 {
			assert.Equal(t, createdConnection.ID, connections[0].ID)
			assert.Equal(t, "Potential_Partner", connections[0].ConnectionType)
			assert.Equal(t, biz.ID, connections[0].TargetEntityID)
		}
	})

	t.Run("Get Connections For Different Source (Empty)", func(t *testing.T) {
		// Build path safely from constants
		getSourcePath := strings.Replace(constants.AppRoutes.InferredBySource, ":entityType", "business", 1)
		getSourcePath = strings.Replace(getSourcePath, ":entityID", fmt.Sprintf("%d", biz.ID), 1)
		url := constInferredBase + getSourcePath

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var connections []ports.InferredConnectionResponse
		json.Unmarshal(w.Body.Bytes(), &connections)

		assert.Len(t, connections, 0)
	})
}
