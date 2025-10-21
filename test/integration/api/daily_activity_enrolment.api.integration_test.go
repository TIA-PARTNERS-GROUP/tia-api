package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestDailyActivityEnrolmentAPI_Integration_Lifecycle(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	userA, tokenA := CreateTestUserAndLogin(t, router, "enroluserA@test.com", "ValidPass123!")
	userB, tokenB := CreateTestUserAndLogin(t, router, "enroluserB@test.com", "ValidPass123!")

	activity := models.DailyActivity{Name: "Yoga Session", Description: "Morning yoga"}
	testutil.TestDB.Create(&activity)
	assert.NotZero(t, activity.ID)

	// --- DEFINE ALL PATHS USING CONSTANTS ---

	// Path 1: /api/v1/daily-activities/:id/enrolments
	actBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.DailyActBase
	actSubPath := strings.Replace(constants.AppRoutes.DailyActEnrol, ":id", fmt.Sprintf("%d", activity.ID), 1)
	activityEnrolmentURL := actBasePath + actSubPath

	// Path 2: /api/v1/users/:id/enrolments (for User A)
	userBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase
	userASubPath := strings.Replace(constants.AppRoutes.UserEnrolments, ":id", fmt.Sprintf("%d", userA.ID), 1)
	userAEnrolmentsURL := userBasePath + userASubPath

	// --- TEST FLOW ---

	// 1. User A enrols
	req, _ := http.NewRequest(http.MethodPost, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "User A should be able to enrol - got response: %s", w.Body.String())

	// 2. User B enrols
	req, _ = http.NewRequest(http.MethodPost, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "User B should be able to enrol - got response: %s", w.Body.String())

	// 3. Get enrolments for the ACTIVITY (should be 2)
	req, _ = http.NewRequest(http.MethodGet, activityEnrolmentURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get activity enrolments - got response: %s", w.Body.String())

	var activityEnrolments []models.DailyActivityEnrolment
	if w.Code == http.StatusOK {
		json.Unmarshal(w.Body.Bytes(), &activityEnrolments)
	}
	assert.Len(t, activityEnrolments, 2, "Should be two enrolments for the activity")

	// 4. Get enrolments for USER A (should be 1)
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments - got response: %s", w.Body.String())

	var userAEnrols []models.DailyActivityEnrolment // Use a new variable
	if w.Code == http.StatusOK {
		json.Unmarshal(w.Body.Bytes(), &userAEnrols)
	}

	if assert.Len(t, userAEnrols, 1, "User A should have one enrolment") {
		assert.Equal(t, activity.ID, userAEnrols[0].DailyActivityID)
	}

	// 5. User A withdraws
	req, _ = http.NewRequest(http.MethodDelete, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, "User A should be able to withdraw")

	// 6. Get enrolments for the ACTIVITY again (should be 1)
	req, _ = http.NewRequest(http.MethodGet, activityEnrolmentURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get updated activity enrolments")

	activityEnrolments = []models.DailyActivityEnrolment{} // Re-init slice
	if w.Code == http.StatusOK {
		json.Unmarshal(w.Body.Bytes(), &activityEnrolments)
	}

	assert.Len(t, activityEnrolments, 1, "Should now be only one enrolment")
	if len(activityEnrolments) > 0 {
		assert.Equal(t, userB.ID, activityEnrolments[0].User.ID, "Remaining user should be User B")
	}

	// 7. Get enrolments for USER A again (should be 0)
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments after withdrawal")

	userAEnrols = []models.DailyActivityEnrolment{} // Re-init slice
	if w.Code == http.StatusOK {
		json.Unmarshal(w.Body.Bytes(), &userAEnrols)
	}
	assert.Len(t, userAEnrols, 0, "User A should have zero enrolments after withdrawing")
}
