package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

// Helper to create a Subscription Plan in the DB
func CreateTestSubscriptionPlan(t *testing.T) models.Subscription {
	plan := models.Subscription{
		Name:        "Test Plan",
		Price:       10.00,
		ValidDays:   IntPtr(30),
		ValidMonths: nil,
	}
	result := testutil.TestDB.Create(&plan)
	assert.NoError(t, result.Error)
	assert.NotZero(t, plan.ID)
	return plan
}

// Helper to directly create a UserSubscription record for setup
func CreateTestUserSubscription(t *testing.T, userID uint, planID uint, isTrial bool) models.UserSubscription {
	userSub := models.UserSubscription{
		UserID:         userID,
		SubscriptionID: planID,
		DateFrom:       time.Now().Add(-24 * time.Hour),     // Start yesterday
		DateTo:         time.Now().Add(30 * 24 * time.Hour), // Active for a month
		IsTrial:        isTrial,
	}
	result := testutil.TestDB.Create(&userSub)
	assert.NoError(t, result.Error)
	assert.NotZero(t, userSub.ID)
	return userSub
}

func TestUserSubscriptionAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// 1. Setup Data
	user, userToken := CreateTestUserAndLogin(t, router, "usub.user@test.com", "ValidPass123!")
	// --- FIX: Capture otherUser model to get its ID ---
	otherUser, otherToken := CreateTestUserAndLogin(t, router, "usub.other@test.com", "ValidPass123!")
	plan := CreateTestSubscriptionPlan(t)

	// Create an active user subscription record
	activeSub := CreateTestUserSubscription(t, user.ID, plan.ID, false)

	// Create an expired user subscription record (should NOT be returned by GET)
	expiredSub := models.UserSubscription{
		UserID:         user.ID,
		SubscriptionID: plan.ID,
		DateFrom:       time.Now().AddDate(-1, 0, 0),
		DateTo:         time.Now().AddDate(-1, 0, 1), // Expired 1 year ago
		IsTrial:        false,
	}
	result := testutil.TestDB.Create(&expiredSub)
	assert.NoError(t, result.Error)

	// Base URLs
	userSubsURL := fmt.Sprintf("%s/users/%d/subscriptions", constants.AppRoutes.APIPrefix, user.ID)
	cancelSubURL := fmt.Sprintf("%s/users/%d/subscriptions/%d", constants.AppRoutes.APIPrefix, user.ID, activeSub.ID)

	t.Run("Get Subscriptions - Forbidden (Other User)", func(t *testing.T) {
		// This uses 'otherToken' to authenticate, but tries to access 'user.ID's subs
		req, _ := http.NewRequest(http.MethodGet, userSubsURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Get Subscriptions - Success (Self)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userSubsURL, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var subs []ports.UserSubscriptionResponse
		json.Unmarshal(w.Body.Bytes(), &subs)
		assert.Equal(t, 1, len(subs), "Should only return one active subscription")
		assert.Equal(t, activeSub.ID, subs[0].ID)
		assert.Equal(t, plan.Name, subs[0].Subscription.Name)
	})

	t.Run("Cancel Subscription - Forbidden (Wrong User ID in URL)", func(t *testing.T) {
		// Auth: 'userToken' (ID 1), trying to cancel subscription for URL ID 'otherUser.ID' (ID 2)
		wrongUserURL := fmt.Sprintf("%s/users/%d/subscriptions/%d", constants.AppRoutes.APIPrefix, otherUser.ID, activeSub.ID)
		req, _ := http.NewRequest(http.MethodDelete, wrongUserURL, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Cancel Subscription - Forbidden (Not Owner)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, cancelSubURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) // Auth is 'otherUser'
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Cancel Subscription - Success (Self)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, cancelSubURL, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Cancellation", func(t *testing.T) {
		// Check that the active subscription list is now empty
		req, _ := http.NewRequest(http.MethodGet, userSubsURL, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var subs []ports.UserSubscriptionResponse
		json.Unmarshal(w.Body.Bytes(), &subs)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 0, len(subs), "The subscription should be deleted (canceled)")
	})
}
