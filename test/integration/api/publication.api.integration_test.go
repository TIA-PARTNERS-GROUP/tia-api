package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/constants"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestPublicationAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	constApiPrefix := constants.AppRoutes.APIPrefix
	constPubBase := constApiPrefix + constants.AppRoutes.PublicationBase

	authorUser, authorToken := CreateTestUserAndLogin(t, router, "pub.author@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "pub.other@test.com", "ValidPass123!")

	var createdPub ports.PublicationResponse

	t.Run("Create Publication", func(t *testing.T) {
		createDTO := ports.CreatePublicationInput{
			UserID:          authorUser.ID,
			PublicationType: models.PublicationArticle,
			Title:           "My First Great Article",
			Content:         "This is the content.",
			Published:       BoolPtr(true),
		}
		body, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, constPubBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authorToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Publication creation failed")
		json.Unmarshal(w.Body.Bytes(), &createdPub)
		assert.NotZero(t, createdPub.ID)
		assert.Equal(t, "my-first-great-article", createdPub.Slug)
		assert.True(t, createdPub.Published)
		assert.NotNil(t, createdPub.PublishedAt)
		createdPubID := createdPub.ID

		req2, _ := http.NewRequest(http.MethodPost, constPubBase, bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+authorToken)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code)

		createDTO.Title = "Other User's Article"
		createDTO.UserID = 9999
		body3, _ := json.Marshal(createDTO)
		req3, _ := http.NewRequest(http.MethodPost, constPubBase, bytes.NewBuffer(body3))
		req3.Header.Set("Content-Type", "application/json")
		req3.Header.Set("Authorization", "Bearer "+authorToken)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusForbidden, w3.Code)

		createdPub.ID = createdPubID 
	})

	t.Run("Get All Publications", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, constPubBase, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var publications []ports.PublicationResponse
		json.Unmarshal(w.Body.Bytes(), &publications)
		assert.Equal(t, 1, len(publications))
		assert.Equal(t, createdPub.ID, publications[0].ID)
	})

	t.Run("Get Publication by ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/id/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedPub ports.PublicationResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedPub)
		assert.Equal(t, createdPub.ID, fetchedPub.ID)
	})

	t.Run("Get Publication by Slug", func(t *testing.T) {
		url := fmt.Sprintf("%s/slug/%s", constPubBase, createdPub.Slug)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedPub ports.PublicationResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedPub)
		assert.Equal(t, createdPub.ID, fetchedPub.ID)
	})

	t.Run("Update Publication - Forbidden (Not Author)", func(t *testing.T) {
		updateDTO := `{"title": "Forbidden Update"}`
		url := fmt.Sprintf("%s/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update Publication - Success (Author)", func(t *testing.T) {
		updatedTitle := "Updated Article Title"
		updateDTO := ports.UpdatePublicationInput{Title: &updatedTitle}
		body, _ := json.Marshal(updateDTO)
		url := fmt.Sprintf("%s/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authorToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedPub ports.PublicationResponse
		json.Unmarshal(w.Body.Bytes(), &updatedPub)
		assert.Equal(t, updatedTitle, updatedPub.Title)
		assert.Equal(t, "updated-article-title", updatedPub.Slug) 
	})

	t.Run("Delete Publication - Forbidden (Not Author)", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Delete Publication - Success (Author)", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+authorToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		url := fmt.Sprintf("%s/id/%d", constPubBase, createdPub.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+authorToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
