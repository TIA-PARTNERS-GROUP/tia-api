package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants" // <-- IMPORT
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestFeedbackAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// --- USE CONSTANTS ---
	constApiPrefix := constants.AppRoutes.APIPrefix
	constFeedbackBase := constApiPrefix + constants.AppRoutes.FeedbackBase

	// 1. Setup: Create an admin user and log in to get a token
	// --- USE THE HELPER ---
	_, token := CreateTestUserAndLogin(t, router, "feedbackadmin@test.com", "ValidPass123!")

	var createdFeedback models.Feedback

	t.Run("Create Feedback (Public)", func(t *testing.T) {
		createDTO := ports.CreateFeedbackInput{
			Name:    "Test User",
			Email:   "test@user.com",
			Content: "This is a test feedback message.",
		}
		body, _ := json.Marshal(createDTO)

		// Use constant path
		req, _ := http.NewRequest(http.MethodPost, constFeedbackBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		json.Unmarshal(w.Body.Bytes(), &createdFeedback)
		assert.NotZero(t, createdFeedback.ID)
		assert.Equal(t, "Test User", createdFeedback.Name)
		assert.Equal(t, "This is a test feedback message.", createdFeedback.Content)
	})

	t.Run("Fail to Get All Feedback (No Auth)", func(t *testing.T) {
		// Use constant path
		req, _ := http.NewRequest(http.MethodGet, constFeedbackBase, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		// Expect 401 Unauthorized because this endpoint is protected
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get All Feedback (Admin Auth)", func(t *testing.T) {
		// Use constant path
		req, _ := http.NewRequest(http.MethodGet, constFeedbackBase, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var feedbacks []models.Feedback
		json.Unmarshal(w.Body.Bytes(), &feedbacks)
		assert.Len(t, feedbacks, 1)
		assert.Equal(t, createdFeedback.ID, feedbacks[0].ID)
	})

	t.Run("Get Feedback By ID (Admin Auth)", func(t *testing.T) {
		// Use constant path
		url := fmt.Sprintf("%s/%d", constFeedbackBase, createdFeedback.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var feedback models.Feedback
		json.Unmarshal(w.Body.Bytes(), &feedback)
		assert.Equal(t, createdFeedback.ID, feedback.ID)
	})

	t.Run("Delete Feedback (Admin Auth)", func(t *testing.T) {
		// Use constant path
		url := fmt.Sprintf("%s/%d", constFeedbackBase, createdFeedback.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		// Use constant path
		url := fmt.Sprintf("%s/%d", constFeedbackBase, createdFeedback.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Expect 404 Not Found because it was deleted
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
