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
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports" // <-- IMPORT PORTS
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

	// Paths using constants
	actBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.DailyActBase
	actSubPath := strings.Replace(constants.AppRoutes.DailyActEnrol, ":id", fmt.Sprintf("%d", activity.ID), 1)
	activityEnrolmentURL := actBasePath + actSubPath

	userBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase
	userASubPath := strings.Replace(constants.AppRoutes.UserEnrolments, ":id", fmt.Sprintf("%d", userA.ID), 1)
	userAEnrolmentsURL := userBasePath + userASubPath

	// --- Test Flow ---

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

	// --- FIX: Unmarshal into ActivityEnrolmentResponse DTO ---
	var activityEnrolmentsDTO []ports.ActivityEnrolmentResponse
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &activityEnrolmentsDTO)
		assert.NoError(t, err)
	}
	assert.Len(t, activityEnrolmentsDTO, 2, "Should be two enrolments for the activity")

	// 4. Get enrolments for USER A (should be 1)
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments - got response: %s", w.Body.String())

	// --- FIX: Unmarshal into UserEnrolmentResponse DTO ---
	var userAEnrolsDTO []ports.UserEnrolmentResponse // Use DTO type
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &userAEnrolsDTO)
		assert.NoError(t, err)
	}

	// --- FIX: Assert against fields in the DTO ---
	if assert.Len(t, userAEnrolsDTO, 1, "User A should have one enrolment") {
		// Assert based on the UserEnrolmentResponse structure
		assert.Equal(t, userA.ID, userAEnrolsDTO[0].UserID)
		assert.Equal(t, activity.ID, userAEnrolsDTO[0].DailyActivity.ID)     // Check ID within nested struct
		assert.Equal(t, activity.Name, userAEnrolsDTO[0].DailyActivity.Name) // Check name within nested struct
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

	activityEnrolmentsDTO = []ports.ActivityEnrolmentResponse{} // Re-init slice
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &activityEnrolmentsDTO)
		assert.NoError(t, err)
	}

	assert.Len(t, activityEnrolmentsDTO, 1, "Should now be only one enrolment")
	if len(activityEnrolmentsDTO) > 0 {
		assert.Equal(t, userB.ID, activityEnrolmentsDTO[0].User.ID, "Remaining user should be User B") // Check ID within nested struct
	}

	// 7. Get enrolments for USER A again (should be 0)
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments after withdrawal")

	userAEnrolsDTO = []ports.UserEnrolmentResponse{} // Re-init slice
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &userAEnrolsDTO)
		assert.NoError(t, err)
	}
	assert.Len(t, userAEnrolsDTO, 0, "User A should have zero enrolments after withdrawing")
}
