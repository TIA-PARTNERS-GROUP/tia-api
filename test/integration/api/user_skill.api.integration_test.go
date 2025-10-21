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

func TestUserSkillAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	user, userToken := CreateTestUserAndLogin(t, router, "uskill.user@test.com", "ValidPass123!")
	otherUser, _ := CreateTestUserAndLogin(t, router, "uskill.other@test.com", "ValidPass123!")

	skill1 := CreateTestSkill(t, "Rust", 1)
	skill2 := CreateTestSkill(t, "DevOps", 2)

	userSkillBaseURL := fmt.Sprintf("%s/users/%d/skills", constants.AppRoutes.APIPrefix, user.ID)
	otherUserSkillBaseURL := fmt.Sprintf("%s/users/%d/skills", constants.AppRoutes.APIPrefix, otherUser.ID)

	t.Run("Add Skill - Success", func(t *testing.T) {
		addDTO := ports.CreateUserSkillInput{
			UserID:           999, 
			SkillID:          skill1.ID,
			ProficiencyLevel: models.ProficiencyIntermediate,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, userSkillBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "User skill creation failed")
		var resp ports.UserSkillResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, skill1.ID, resp.SkillID)
		assert.Equal(t, user.ID, resp.UserID)
		assert.Equal(t, models.ProficiencyIntermediate, resp.ProficiencyLevel)
		assert.Equal(t, "Rust", resp.Skill.Name)
	})

	t.Run("Add Skill - Forbidden (Targeting Another User)", func(t *testing.T) {
		addDTO := ports.CreateUserSkillInput{
			SkillID:          skill2.ID,
			ProficiencyLevel: models.ProficiencyAdvanced,
		}
		body, _ := json.Marshal(addDTO)
		req, _ := http.NewRequest(http.MethodPost, otherUserSkillBaseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken) 
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Get Skills - Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, userSkillBaseURL, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp ports.UserSkillsResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 1, resp.Count)
		assert.Equal(t, skill1.ID, resp.Skills[0].SkillID)
	})

	t.Run("Update Skill - Forbidden (Targeting Another User)", func(t *testing.T) {
		updateDTO := `{"proficiency_level": "expert"}`
		url := fmt.Sprintf("%s/%d", otherUserSkillBaseURL, skill1.ID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(updateDTO)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update Skill - Success", func(t *testing.T) {
		updatedProficiency := models.ProficiencyAdvanced
		updateDTO := ports.UpdateUserSkillInput{ProficiencyLevel: &updatedProficiency}
		body, _ := json.Marshal(updateDTO)
		url := fmt.Sprintf("%s/%d", userSkillBaseURL, skill1.ID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedSkill ports.UserSkillResponse
		json.Unmarshal(w.Body.Bytes(), &updatedSkill)
		assert.Equal(t, models.ProficiencyAdvanced, updatedSkill.ProficiencyLevel)
	})

	t.Run("Remove Skill - Forbidden (Targeting Another User)", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", otherUserSkillBaseURL, skill1.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Remove Skill - Success", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", userSkillBaseURL, skill1.ID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Removal", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", userSkillBaseURL, skill1.ID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
