package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" 
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" 
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)
func TestL2EAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	
	user, token := CreateTestUserAndLogin(t, router, "l2e-user@test.com", "ValidPass123!")
	
	constL2EBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.L2EBase
	constUserBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.UsersBase
	var createdResponse ports.L2EResponseResponse
	t.Run("Create L2E Response", func(t *testing.T) {
		
		responsePayload := `{"question1": "answer1", "score": 10}`
		responseJSON := datatypes.JSON(responsePayload)
		createDTO := ports.CreateL2EResponseInput{
			Response: responseJSON,
			
		}
		body, _ := json.Marshal(createDTO)
		
		req, _ := http.NewRequest(http.MethodPost, constL2EBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		err := json.Unmarshal(w.Body.Bytes(), &createdResponse)
		assert.NoError(t, err, "Failed to unmarshal response")
		assert.NotZero(t, createdResponse.ID)
		assert.Equal(t, user.ID, createdResponse.UserID)
		
		assert.JSONEq(t, responsePayload, string(createdResponse.Response))
	})
	t.Run("Get L2E Responses For User", func(t *testing.T) {
		
		userL2ESubPath := strings.Replace(constants.AppRoutes.UserL2EResponses, ":id", fmt.Sprintf("%d", user.ID), 1)
		url := constUserBase + userL2ESubPath
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var responses []ports.L2EResponseResponse
		err := json.Unmarshal(w.Body.Bytes(), &responses)
		assert.NoError(t, err, "Failed to unmarshal responses")
		assert.Len(t, responses, 1)
		if len(responses) > 0 {
			assert.Equal(t, createdResponse.ID, responses[0].ID)
			assert.Equal(t, user.ID, responses[0].UserID)
			assert.JSONEq(t, string(createdResponse.Response), string(responses[0].Response))
		}
	})
	t.Run("Get L2E Responses For Other User (Empty)", func(t *testing.T) {
		nonExistentUserID := uint(9999)
		
		userL2ESubPath := strings.Replace(constants.AppRoutes.UserL2EResponses, ":id", fmt.Sprintf("%d", nonExistentUserID), 1)
		url := constUserBase + userL2ESubPath
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code) 
		var responses []ports.L2EResponseResponse
		err := json.Unmarshal(w.Body.Bytes(), &responses)
		assert.NoError(t, err, "Failed to unmarshal responses")
		assert.Len(t, responses, 0)
	})
}
