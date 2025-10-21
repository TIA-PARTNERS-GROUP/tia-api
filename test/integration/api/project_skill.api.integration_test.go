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

// Helper to create a Skill in the DB
func CreateTestSkill(t *testing.T, name string, id uint) models.Skill {
	skill := models.Skill{
		ID:       id,
		Category: "Test Category",
		Name:     name,
		Active:   true,
	}
	result := testutil.TestDB.FirstOrCreate(&skill)
	assert.NoError(t, result.Error, "Failed to create test skill")
	return skill
}

func TestProjectSkillAPI_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	// 1. Create Users
	managerUser, managerToken := CreateTestUserAndLogin(t, router, "skill.manager@test.com", "ValidPass123!")
	_, otherToken := CreateTestUserAndLogin(t, router, "skill.other@test.com", "ValidPass123!")

	// 2. Create Project and Skills
	project := CreateTestProjectHelper(t, router, managerUser, managerToken)
	projectID := project.ID
	skill1 := CreateTestSkill(t, "GoLang", 101)
	skill2 := CreateTestSkill(t, "AWS", 102)

	skillsBaseURL := fmt.Sprintf("%s/projects/%d/skills", constants.AppRoutes.APIPrefix, projectID)
	skill1SpecificURL := fmt.Sprintf("%s/%d", skillsBaseURL, skill1.ID)

	t.Run("Add Skill - Forbidden (Not Manager)", func(t *testing.T) {
		addDTO := ports.CreateProjectSkillInput{
			SkillID:    skill1.ID,
			Importance: models.SkillImportanceRequired,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, skillsBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Add Skill 1 - Success (Manager)", func(t *testing.T) {
		addDTO := ports.CreateProjectSkillInput{
			ProjectID:  999, // Should be ignored
			SkillID:    skill1.ID,
			Importance: models.SkillImportanceRequired,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, skillsBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var newSkill ports.ProjectSkillResponse
		json.Unmarshal(w.Body.Bytes(), &newSkill)
		assert.Equal(t, skill1.ID, newSkill.SkillID)
		assert.Equal(t, models.SkillImportanceRequired, newSkill.Importance)
	})

	t.Run("Add Skill 2 - Success (Manager)", func(t *testing.T) {
		addDTO := ports.CreateProjectSkillInput{
			SkillID:    skill2.ID,
			Importance: models.SkillImportancePreferred,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, skillsBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Get Project Skills", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, skillsBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken) // Any auth user
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.ProjectSkillsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 2, resp.Count)

		// --- MODIFIED ASSERTIONS FOR ORDERING ---
		// Find the two skills by ID to assert their properties, regardless of index.
		var foundSkill1, foundSkill2 *ports.ProjectSkillResponse
		for i := range resp.Skills {
			if resp.Skills[i].SkillID == skill1.ID {
				foundSkill1 = &resp.Skills[i]
			}
			if resp.Skills[i].SkillID == skill2.ID {
				foundSkill2 = &resp.Skills[i]
			}
		}

		assert.NotNil(t, foundSkill1, "Skill 1 (GoLang) not found")
		assert.NotNil(t, foundSkill2, "Skill 2 (AWS) not found")

		// Assert properties, not index/order
		assert.Equal(t, models.SkillImportanceRequired, foundSkill1.Importance, "Skill 1 importance incorrect")
		assert.Equal(t, models.SkillImportancePreferred, foundSkill2.Importance, "Skill 2 importance incorrect")
		// --- END MODIFIED ASSERTIONS ---
	})

	t.Run("Update Skill - Forbidden (Not Manager)", func(t *testing.T) {
		updateDTO := `{"importance": "optional"}`
		req, _ := http.NewRequest(http.MethodPut, skill1SpecificURL, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update Skill - Success (Manager)", func(t *testing.T) {
		newImportance := models.SkillImportanceOptional
		updateDTO := ports.UpdateProjectSkillInput{Importance: &newImportance}
		body, _ := json.Marshal(updateDTO)
		req, _ := http.NewRequest(http.MethodPut, skill1SpecificURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedSkill ports.ProjectSkillResponse
		json.Unmarshal(w.Body.Bytes(), &updatedSkill)
		assert.Equal(t, models.SkillImportanceOptional, updatedSkill.Importance)
	})

	t.Run("Remove Skill - Forbidden (Not Manager)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, skill1SpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Remove Skill - Success (Manager)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, skill1SpecificURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Removal", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, skillsBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+managerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.ProjectSkillsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp.Count)
		assert.Equal(t, skill2.ID, resp.Skills[0].SkillID)
	})
}
