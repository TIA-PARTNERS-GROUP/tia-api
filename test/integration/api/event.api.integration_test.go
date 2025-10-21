package main
import (
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
func TestEventAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()
	user, token := CreateTestUserAndLogin(t, router, "eventuser@test.com", "ValidPass123!")
	constEventBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.EventBase
	eventInput := map[string]interface{}{
		"event_type": "user_action",
		"payload": map[string]interface{}{
			"action": "test_action",
			"data":   "test_data",
		},
	}
	eventInputPayloadJSON, _ := json.Marshal(eventInput["payload"]) 
	req, _ := http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Should create event - got response: %s", w.Body.String())
	
	var createdEventResponse ports.EventResponse
	err := json.Unmarshal(w.Body.Bytes(), &createdEventResponse)
	assert.NoError(t, err, "Failed to unmarshal create event response")
	
	
	assert.NotZero(t, createdEventResponse.ID)
	assert.Equal(t, "user_action", createdEventResponse.EventType)                         
	if assert.NotNil(t, createdEventResponse.UserID, "UserID pointer should not be nil") { 
		assert.Equal(t, user.ID, *createdEventResponse.UserID, "Event should be associated with the authenticated user") 
	}
	assert.JSONEq(t, string(eventInputPayloadJSON), string(createdEventResponse.Payload), "Payload does not match") 
	
	
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", constEventBase, createdEventResponse.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get event by ID - got response: %s", w.Body.String())
	
	var fetchedEventResponse ports.EventResponse
	err = json.Unmarshal(w.Body.Bytes(), &fetchedEventResponse)
	assert.NoError(t, err, "Failed to unmarshal get event by ID response")
	
	
	assert.Equal(t, createdEventResponse.ID, fetchedEventResponse.ID)
	assert.Equal(t, createdEventResponse.EventType, fetchedEventResponse.EventType)
	assert.Equal(t, createdEventResponse.UserID, fetchedEventResponse.UserID) 
	assert.JSONEq(t, string(createdEventResponse.Payload), string(fetchedEventResponse.Payload))
	
	
	req, _ = http.NewRequest(http.MethodGet, constEventBase, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get all events - got response: %s", w.Body.String())
	
	var eventsResponse []ports.EventResponse
	err = json.Unmarshal(w.Body.Bytes(), &eventsResponse)
	assert.NoError(t, err, "Failed to unmarshal get all events response")
	
	if assert.Len(t, eventsResponse, 1) {
		assert.Equal(t, createdEventResponse.ID, eventsResponse[0].ID) 
	}
	
	
	
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"?event_type=user_action", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get filtered events - got response: %s", w.Body.String())
	
	eventsResponse = []ports.EventResponse{} 
	err = json.Unmarshal(w.Body.Bytes(), &eventsResponse)
	assert.NoError(t, err, "Failed to unmarshal get filtered events response")
	
	if assert.Len(t, eventsResponse, 1) {
		assert.Equal(t, createdEventResponse.ID, eventsResponse[0].ID) 
	}
	
	
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Content-Type", "application/json") 
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should reject unauthorized event creation")
	
	invalidEventInput := map[string]interface{}{"event_type": "", "payload": eventInput["payload"]}
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code) 
	
	invalidEventInput2 := map[string]interface{}{"event_type": "test_event"} 
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput2))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code) 
	
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent event")
	
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject invalid event ID format")
}
