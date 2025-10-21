package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestBusinessAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// --- USE THE HELPER ---
	user, token := CreateTestUserAndLogin(t, router, "bizowner@api.com", "ValidPassword123!")

	// --- USE CONSTANTS ---
	constApiBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.BusinessBase

	createDTO := ports.CreateBusinessInput{
		OperatorUserID:   user.ID,
		Name:             "API Test Corp",
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseGrowth,
	}
	body, _ := json.Marshal(createDTO)

	// Use constant path
	req, _ := http.NewRequest(http.MethodPost, constApiBase, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdBusiness ports.BusinessResponse
	json.Unmarshal(w.Body.Bytes(), &createdBusiness)
	assert.Equal(t, "API Test Corp", createdBusiness.Name)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", constApiBase, createdBusiness.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Use constant path
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", constApiBase, createdBusiness.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", constApiBase, createdBusiness.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
