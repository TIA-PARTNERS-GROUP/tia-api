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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func CreateTestProject(t *testing.T, router *gin.Engine, managerUser models.User, managerToken string) ports.ProjectResponse {
	constProjectBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.ProjectBase
	createDTO := ports.CreateProjectInput{
		ManagedByUserID: managerUser.ID,
		Name:            "Applicant Test Project",
		ProjectStatus:   models.ProjectStatusActive,
	}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, constProjectBase, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+managerToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Project creation for test failed")
	var createdProject ports.ProjectResponse
	json.Unmarshal(w.Body.Bytes(), &createdProject)
	return createdProject
}

func TestProjectApplicantAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	managerUser, managerToken := CreateTestUserAndLogin(t, router, "app.manager@test.com", "ValidPass123!")
	applicantUser, applicantToken := CreateTestUserAndLogin(t, router, "app.applicant@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "app.other@test.com", "ValidPass123!")

	project := CreateTestProject(t, router, managerUser, managerToken)
	projectID := project.ID
	applicantID := applicantUser.ID

	applyURL := fmt.Sprintf("%s/projects/%d/apply", constants.AppRoutes.APIPrefix, projectID)
	applicantsURL := fmt.Sprintf("%s/projects/%d/applicants", constants.AppRoutes.APIPrefix, projectID)
	myApplicationsURL := fmt.Sprintf("%s/users/%d/applications", constants.AppRoutes.APIPrefix, applicantID)

	t.Run("Applicant applies to project", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, applyURL, nil)
		req.Header.Set("Authorization", "Bearer "+applicantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Applicant applies again (fail)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, applyURL, nil)
		req.Header.Set("Authorization", "Bearer "+applicantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusConflict, w.Code) 
	})

	t.Run("Manager gets applicants", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, applicantsURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var applicants []ports.ProjectApplicantResponse
		json.Unmarshal(w.Body.Bytes(), &applicants)
		assert.Equal(t, 1, len(applicants))
		assert.Equal(t, applicantID, applicants[0].UserID)
	})

	t.Run("Other user gets applicants (fail)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, applicantsURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Applicant gets their applications", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, myApplicationsURL, nil)
		req.Header.Set("Authorization", "Bearer "+applicantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var applications []ports.UserApplicationResponse
		json.Unmarshal(w.Body.Bytes(), &applications)
		assert.Equal(t, 1, len(applications))
		assert.Equal(t, projectID, applications[0].ProjectID)
		assert.Equal(t, project.Name, applications[0].Project.Name)
	})

	t.Run("Other user gets applicant's applications (fail)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, myApplicationsURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Manager gets applicant's applications (fail)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, myApplicationsURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Applicant withdraws application", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, applyURL, nil)
		req.Header.Set("Authorization", "Bearer "+applicantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Applicant withdraws again (fail)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, applyURL, nil)
		req.Header.Set("Authorization", "Bearer "+applicantToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code) 
	})

	t.Run("Manager gets applicants (empty)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, applicantsURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var applicants []ports.ProjectApplicantResponse
		json.Unmarshal(w.Body.Bytes(), &applicants)
		assert.Equal(t, 0, len(applicants))
	})
}
