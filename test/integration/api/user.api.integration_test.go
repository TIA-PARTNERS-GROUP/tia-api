package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestUserAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	createDTO := ports.UserCreationSchema{
		FirstName:  "ApiTest",
		LastName:   "User",
		LoginEmail: "api@test.com",
		Password:   "ValidPassword123!",
	}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser ports.UserResponse
	json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.Equal(t, "ApiTest", createdUser.FirstName)
	assert.NotZero(t, createdUser.ID)

	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedUser ports.UserResponse
	json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)

	updateDTO := `{"first_name": "ApiUpdated"}`
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%d", createdUser.ID), bytes.NewBuffer([]byte(updateDTO)))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedUser ports.UserResponse
	json.Unmarshal(w.Body.Bytes(), &updatedUser)
	assert.Equal(t, "ApiUpdated", updatedUser.FirstName)

	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
