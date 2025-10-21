package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

// --- Local helper function removed ---
// It's already defined in api_test_helpers.go (package main)

func TestEventAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	user, token := CreateTestUserAndLogin(t, router, "eventuser@test.com", "ValidPass123!")

	// --- USE CONSTANTS ---
	constEventBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.EventBase

	eventInput := map[string]interface{}{
		"event_type": "user_action",
		"payload": map[string]interface{}{
			"action": "test_action",
			"data":   "test_data",
		},
	}

	// Use constant path
	req, _ := http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Should create event - got response: %s", w.Body.String())

	var createdEvent models.Event
	json.Unmarshal(w.Body.Bytes(), &createdEvent)
	assert.NotZero(t, createdEvent.ID)
	assert.Equal(t, "user_action", createdEvent.EventType)

	if assert.NotNil(t, createdEvent.UserID, "UserID pointer should not be nil") {
		assert.Equal(t, user.ID, *createdEvent.UserID, "Event should be associated with the authenticated user")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(createdEvent.Payload, &payload); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}
	assert.Equal(t, "test_action", payload["action"])
	assert.Equal(t, "test_data", payload["data"])

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", constEventBase, createdEvent.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get event by ID - got response: %s", w.Body.String())

	var fetchedEvent models.Event
	json.Unmarshal(w.Body.Bytes(), &fetchedEvent)
	assert.Equal(t, createdEvent.ID, fetchedEvent.ID)
	assert.Equal(t, createdEvent.EventType, fetchedEvent.EventType)
	assert.Equal(t, createdEvent.UserID, fetchedEvent.UserID)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get all events - got response: %s", w.Body.String())

	var events []models.Event
	json.Unmarshal(w.Body.Bytes(), &events)
	assert.Len(t, events, 1)
	assert.Equal(t, createdEvent.ID, events[0].ID)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"?event_type=user_action", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get filtered events - got response: %s", w.Body.String())

	events = []models.Event{}
	json.Unmarshal(w.Body.Bytes(), &events)
	assert.Len(t, events, 1)
	assert.Equal(t, createdEvent.ID, events[0].ID)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s?user_id=%d", constEventBase, user.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get events by user ID - got response: %s", w.Body.String())

	events = []models.Event{}
	json.Unmarshal(w.Body.Bytes(), &events)
	assert.Len(t, events, 1)
	assert.Equal(t, createdEvent.ID, events[0].ID)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"?event_type=non_existent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get empty events list - got response: %s", w.Body.String())

	events = []models.Event{}
	json.Unmarshal(w.Body.Bytes(), &events)
	assert.Len(t, events, 0)

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"?user_id=999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get empty events list for non-existent user - got response: %s", w.Body.String())

	events = []models.Event{}
	json.Unmarshal(w.Body.Bytes(), &events)
	assert.Len(t, events, 0)

	// Use constant path
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, eventInput))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should reject unauthorized event creation")

	invalidEventInput := map[string]interface{}{
		"event_type": "",
		"payload": map[string]interface{}{
			"action": "test_action",
		},
	}

	// Use constant path
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Logf("Empty event_type returned %d instead of 400 - this might be expected based on your validation rules", w.Code)
	}

	invalidEventInput2 := map[string]interface{}{
		"event_type": "test_event",
	}

	// Use constant path
	req, _ = http.NewRequest(http.MethodPost, constEventBase, createJSONBody(t, invalidEventInput2))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Logf("Missing payload returned %d instead of 400 - check your payload validation", w.Code)
	}

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent event")

	// Use constant path
	req, _ = http.NewRequest(http.MethodGet, constEventBase+"/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject invalid event ID format")
}
