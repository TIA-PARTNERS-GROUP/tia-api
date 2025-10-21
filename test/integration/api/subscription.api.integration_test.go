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

func TestSubscriptionAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	constApiPrefix := constants.AppRoutes.APIPrefix
	constSubBase := constApiPrefix + constants.AppRoutes.SubscriptionBase

	// 1. Create Users
	_, userToken := CreateTestUserAndLogin(t, router, "sub.user@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "sub.other@test.com", "ValidPass123!")

	var createdPlan ports.SubscriptionResponse

	t.Run("Create Subscription Plan", func(t *testing.T) {
		createDTO := ports.CreateSubscriptionInput{
			Name:        "Pro Annual",
			Price:       99.99,
			ValidMonths: IntPtr(12),
		}
		body, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, constSubBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Subscription plan creation failed")
		json.Unmarshal(w.Body.Bytes(), &createdPlan)
		assert.NotZero(t, createdPlan.ID)
		assert.Equal(t, float64(12), float64(*createdPlan.ValidMonths))

		// Attempt to create duplicate (should fail)
		req2, _ := http.NewRequest(http.MethodPost, constSubBase, bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+userToken)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code)
	})

	t.Run("Get Subscription Plan by ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constSubBase, createdPlan.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedPlan ports.SubscriptionResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedPlan)
		assert.Equal(t, createdPlan.ID, fetchedPlan.ID)
	})

}
