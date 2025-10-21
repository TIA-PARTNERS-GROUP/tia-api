package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // Keep models for checking against DB if needed
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"     // <-- IMPORT PORTS
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	// Import datatypes for payload check
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
	eventInputPayloadJSON, _ := json.Marshal(eventInput["payload"]) // For later comparison

	req, _ := http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Should create event - got response: %s", w.Body.String())

	// --- FIX: Unmarshal into ports.EventResponse DTO ---
	var createdEventResponse ports.EventResponse
	err := json.Unmarshal(w.Body.Bytes(), &createdEventResponse)
	assert.NoError(t, err, "Failed to unmarshal create event response")
	// --- End Fix ---

	// --- FIX: Assert against fields in the DTO ---
	assert.NotZero(t, createdEventResponse.ID)
	assert.Equal(t, "user_action", createdEventResponse.EventType)                         // Check DTO field
	if assert.NotNil(t, createdEventResponse.UserID, "UserID pointer should not be nil") { // Check DTO field
		assert.Equal(t, user.ID, *createdEventResponse.UserID, "Event should be associated with the authenticated user") // Check DTO field
	}
	assert.JSONEq(t, string(eventInputPayloadJSON), string(createdEventResponse.Payload), "Payload does not match") // Compare payload JSON
	// --- End Fix ---

	// Get Event By ID - Check response against DTO
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", constEventBase, createdEventResponse.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get event by ID - got response: %s", w.Body.String())

	// --- FIX: Unmarshal into ports.EventResponse DTO ---
	var fetchedEventResponse ports.EventResponse
	err = json.Unmarshal(w.Body.Bytes(), &fetchedEventResponse)
	assert.NoError(t, err, "Failed to unmarshal get event by ID response")
	// --- End Fix ---

	// --- FIX: Assert against fields in the DTO ---
	assert.Equal(t, createdEventResponse.ID, fetchedEventResponse.ID)
	assert.Equal(t, createdEventResponse.EventType, fetchedEventResponse.EventType)
	assert.Equal(t, createdEventResponse.UserID, fetchedEventResponse.UserID) // Compare pointers
	assert.JSONEq(t, string(createdEventResponse.Payload), string(fetchedEventResponse.Payload))
	// --- End Fix ---

	// Get All Events - Check response against DTO array
	req, _ = http.NewRequest(http.MethodGet, constEventBase, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get all events - got response: %s", w.Body.String())

	// --- FIX: Unmarshal into slice of ports.EventResponse DTO ---
	var eventsResponse []ports.EventResponse
	err = json.Unmarshal(w.Body.Bytes(), &eventsResponse)
	assert.NoError(t, err, "Failed to unmarshal get all events response")
	// --- End Fix ---

	if assert.Len(t, eventsResponse, 1) {
		assert.Equal(t, createdEventResponse.ID, eventsResponse[0].ID) // Check DTO field
	}

	// ... (Rest of the test cases: filtering, unauthorized, bad requests, not found)
	// Make sure any unmarshalling in the remaining cases also uses ports.EventResponse

	// Example: Get filtered events
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"?event_type=user_action", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get filtered events - got response: %s", w.Body.String())

	// --- FIX: Unmarshal into slice of ports.EventResponse DTO ---
	eventsResponse = []ports.EventResponse{} // Reset slice
	err = json.Unmarshal(w.Body.Bytes(), &eventsResponse)
	assert.NoError(t, err, "Failed to unmarshal get filtered events response")
	// --- End Fix ---
	if assert.Len(t, eventsResponse, 1) {
		assert.Equal(t, createdEventResponse.ID, eventsResponse[0].ID) // Check DTO field
	}

	// ... continue checking and fixing other parts of the test similarly ...

	// Test Unauthorized POST
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Content-Type", "application/json") // No Auth header
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should reject unauthorized event creation")

	// Test Invalid Input (Empty EventType)
	invalidEventInput := map[string]interface{}{"event_type": "", "payload": eventInput["payload"]}
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code) // Expect validation error

	// Test Invalid Input (Missing Payload)
	invalidEventInput2 := map[string]interface{}{"event_type": "test_event"} // Payload missing
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput2))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code) // Expect validation error (payload required)

	// Test Get Not Found
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent event")

	// Test Get Invalid ID
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject invalid event ID format")

}
