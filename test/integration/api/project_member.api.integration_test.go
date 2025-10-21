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

func CreateTestProjectHelper(t *testing.T, router *gin.Engine, managerUser models.User, managerToken string) ports.ProjectResponse {
	constProjectBase := constants.AppRoutes.APIPrefix + constants.AppRoutes.ProjectBase
	createDTO := ports.CreateProjectInput{
		ManagedByUserID: managerUser.ID,
		Name:            fmt.Sprintf("Member Test Project %d", managerUser.ID),
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

func TestProjectMemberAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	managerUser, managerToken := CreateTestUserAndLogin(t, router, "member.manager@test.com", "ValidPass123!")
	memberUser, memberToken := CreateTestUserAndLogin(t, router, "member.user@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "member.other@test.com", "ValidPass123!")

	project := CreateTestProjectHelper(t, router, managerUser, managerToken)
	projectID := project.ID
	memberID := memberUser.ID

	membersBaseURL := fmt.Sprintf("%s/projects/%d/members", constants.AppRoutes.APIPrefix, projectID)
	memberSpecificURL := fmt.Sprintf("%s/%d", membersBaseURL, memberID)
	myMembershipsURL := fmt.Sprintf("%s/users/%d/project-memberships", constants.AppRoutes.APIPrefix, memberID)

	t.Run("Add Member - Forbidden (Not Manager)", func(t *testing.T) {
		addDTO := ports.AddProjectMemberInput{
			UserID: memberID,
			Role:   models.ProjectMemberRoleContributor,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, membersBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Add Member - Success (Manager)", func(t *testing.T) {
		addDTO := ports.AddProjectMemberInput{
			ProjectID: 999, 
			UserID:    memberID,
			Role:      models.ProjectMemberRoleContributor,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, membersBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var newMember ports.ProjectMemberResponse
		json.Unmarshal(w.Body.Bytes(), &newMember)
		assert.Equal(t, memberID, newMember.UserID)
		assert.Equal(t, projectID, newMember.ProjectID) 
		assert.Equal(t, models.ProjectMemberRoleContributor, newMember.Role)
	})

	t.Run("Get Project Members", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, membersBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.ProjectMembersResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 2, resp.Count) 
	})

	t.Run("Get Specific Project Member", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, memberSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var member ports.ProjectMemberResponse
		json.Unmarshal(w.Body.Bytes(), &member)
		assert.Equal(t, memberID, member.UserID)
	})

	t.Run("Get My Project Memberships", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, myMembershipsURL, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.ProjectMembersResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp.Count)
		assert.Equal(t, projectID, resp.Members[0].ProjectID)
	})

	t.Run("Get My Project Memberships (Filtered)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, myMembershipsURL+"?role=contributor", nil)
		req.Header.Set("Authorization", "Bearer "+memberToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.ProjectMembersResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp.Count)

		req, _ = http.NewRequest(http.MethodGet, myMembershipsURL+"?role=manager", nil)
		req.Header.Set("Authorization", "Bearer "+memberToken) 
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Count)
	})

	t.Run("Update Member Role - Forbidden (Not Manager)", func(t *testing.T) {
		updateDTO := `{"role": "reviewer"}`
		req, _ := http.NewRequest(http.MethodPut, memberSpecificURL, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update Member Role - Success (Manager)", func(t *testing.T) {
		updateDTO := ports.UpdateProjectMemberRoleInput{Role: models.ProjectMemberRoleReviewer}
		body, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, memberSpecificURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedMember ports.ProjectMemberResponse
		json.Unmarshal(w.Body.Bytes(), &updatedMember)
		assert.Equal(t, models.ProjectMemberRoleReviewer, updatedMember.Role)
	})

	t.Run("Remove Member - Forbidden (Not Manager or Self)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, memberSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Remove Member - Success (Self)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, memberSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+memberToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	addDTO := ports.AddProjectMemberInput{
		UserID: memberID,
		Role:   models.ProjectMemberRoleContributor,
	}
	body, _ := json.Marshal(addDTO)
	reqAdd, _ := http.NewRequest(http.MethodPost, membersBaseURL, bytes.NewBuffer(body))
	reqAdd.Header.Set("Content-Type", "application/json")
	reqAdd.Header.Set("Authorization", "Bearer "+managerToken)
	wAdd := httptest.NewRecorder()
	router.ServeHTTP(wAdd, reqAdd)
	assert.Equal(t, http.StatusCreated, wAdd.Code)

	t.Run("Remove Member - Success (Manager)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, memberSpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Remove Manager - Fail (Service Logic)", func(t *testing.T) {
		managerMemberURL := fmt.Sprintf("%s/%d", membersBaseURL, managerUser.ID)
		req, _ := http.NewRequest(http.MethodDelete, managerMemberURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code) 
	})
}
