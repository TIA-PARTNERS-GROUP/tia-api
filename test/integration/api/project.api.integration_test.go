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

func TestProjectAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	constApiPrefix := constants.AppRoutes.APIPrefix
	constProjectBase := constApiPrefix + constants.AppRoutes.ProjectBase

	managerUser, managerToken := CreateTestUserAndLogin(t, router, "proj.manager@test.com", "ValidPass123!")
	_, memberToken := CreateTestUserAndLogin(t, router, "proj.member@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "proj.other@test.com", "ValidPass123!")

	var createdProjectID uint
	var createdProject ports.ProjectResponse

	t.Run("Create Project", func(t *testing.T) {
		createDTO := ports.CreateProjectInput{
			ManagedByUserID: managerUser.ID,
			Name:            "Test API Project",
			ProjectStatus:   models.ProjectStatusPlanning,
		}
		body, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, constProjectBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Project creation failed")
		json.Unmarshal(w.Body.Bytes(), &createdProject)
		assert.NotZero(t, createdProject.ID)
		assert.Equal(t, "Test API Project", createdProject.Name)
		assert.Equal(t, managerUser.ID, createdProject.Manager.ID)
		createdProjectID = createdProject.ID
	})

	t.Run("Get Project By ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedProject ports.ProjectResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedProject)
		assert.Equal(t, createdProjectID, fetchedProject.ID)
	})

	t.Run("Get All Projects", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, constProjectBase, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var projects []ports.ProjectResponse
		json.Unmarshal(w.Body.Bytes(), &projects)
		assert.True(t, len(projects) >= 1)
	})

	t.Run("Update Project - Forbidden (Not Manager)", func(t *testing.T) {
		updateDTO := `{"name": "Forbidden Update"}`
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+memberToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update Project - Success (Manager)", func(t *testing.T) {
		updatedName := "Updated Project Name"
		updateDTO := ports.UpdateProjectInput{Name: &updatedName}
		body, _ := json.Marshal(updateDTO)
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedProject ports.ProjectResponse
		json.Unmarshal(w.Body.Bytes(), &updatedProject)
		assert.Equal(t, updatedName, updatedProject.Name)
	})

	t.Run("Delete Project - Forbidden (Not Manager)", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Delete Project - Success (Manager)", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constProjectBase, createdProjectID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
