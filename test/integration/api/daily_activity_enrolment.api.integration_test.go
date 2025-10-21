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
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports" 
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
	
	actBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.DailyActBase
	actSubPath := strings.Replace(constants.AppRoutes.DailyActEnrol, ":id", fmt.Sprintf("%d", activity.ID), 1)
	activityEnrolmentURL := actBasePath + actSubPath
	userBasePath := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase
	userASubPath := strings.Replace(constants.AppRoutes.UserEnrolments, ":id", fmt.Sprintf("%d", userA.ID), 1)
	userAEnrolmentsURL := userBasePath + userASubPath
	
	
	req, _ := http.NewRequest(http.MethodPost, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "User A should be able to enrol - got response: %s", w.Body.String())
	
	req, _ = http.NewRequest(http.MethodPost, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "User B should be able to enrol - got response: %s", w.Body.String())
	
	req, _ = http.NewRequest(http.MethodGet, activityEnrolmentURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get activity enrolments - got response: %s", w.Body.String())
	
	var activityEnrolmentsDTO []ports.ActivityEnrolmentResponse
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &activityEnrolmentsDTO)
		assert.NoError(t, err)
	}
	assert.Len(t, activityEnrolmentsDTO, 2, "Should be two enrolments for the activity")
	
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments - got response: %s", w.Body.String())
	
	var userAEnrolsDTO []ports.UserEnrolmentResponse 
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &userAEnrolsDTO)
		assert.NoError(t, err)
	}
	
	if assert.Len(t, userAEnrolsDTO, 1, "User A should have one enrolment") {
		
		assert.Equal(t, userA.ID, userAEnrolsDTO[0].UserID)
		assert.Equal(t, activity.ID, userAEnrolsDTO[0].DailyActivity.ID)     
		assert.Equal(t, activity.Name, userAEnrolsDTO[0].DailyActivity.Name) 
	}
	
	req, _ = http.NewRequest(http.MethodDelete, activityEnrolmentURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code, "User A should be able to withdraw")
	
	req, _ = http.NewRequest(http.MethodGet, activityEnrolmentURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get updated activity enrolments")
	activityEnrolmentsDTO = []ports.ActivityEnrolmentResponse{} 
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &activityEnrolmentsDTO)
		assert.NoError(t, err)
	}
	assert.Len(t, activityEnrolmentsDTO, 1, "Should now be only one enrolment")
	if len(activityEnrolmentsDTO) > 0 {
		assert.Equal(t, userB.ID, activityEnrolmentsDTO[0].User.ID, "Remaining user should be User B") 
	}
	
	req, _ = http.NewRequest(http.MethodGet, userAEnrolmentsURL, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should be able to get user A's enrolments after withdrawal")
	userAEnrolsDTO = []ports.UserEnrolmentResponse{} 
	if w.Code == http.StatusOK {
		err := json.Unmarshal(w.Body.Bytes(), &userAEnrolsDTO)
		assert.NoError(t, err)
	}
	assert.Len(t, userAEnrolsDTO, 0, "User A should have zero enrolments after withdrawing")
}
