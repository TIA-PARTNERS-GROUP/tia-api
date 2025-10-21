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

func TestSkillAPI_Integration_CRUD(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	router := SetupRouter()

	constApiPrefix := constants.AppRoutes.APIPrefix
	constSkillBase := constApiPrefix + constants.AppRoutes.SkillsBase

	authorUser, userToken := CreateTestUserAndLogin(t, router, "skill.user@test.com", "ValidPass123!")

	var createdSkill ports.SkillResponse
	var createdSkillID uint

	t.Run("Create Skill", func(t *testing.T) {
		createDTO := ports.CreateSkillInput{
			Category: "Programming",
			Name:     "GoLang",
			Active:   BoolPtr(true),
		}
		body, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, constSkillBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Skill creation failed")
		json.Unmarshal(w.Body.Bytes(), &createdSkill)
		assert.NotZero(t, createdSkill.ID)
		assert.Equal(t, "GoLang", createdSkill.Name)
		createdSkillID = createdSkill.ID

		req2, _ := http.NewRequest(http.MethodPost, constSkillBase, bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("Authorization", "Bearer "+userToken)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code, "Duplicate skill creation did not fail")
	})

	t.Run("Get Skill By ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constSkillBase, createdSkillID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedSkill ports.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &fetchedSkill)
		assert.Equal(t, createdSkillID, fetchedSkill.ID)
	})

	t.Run("Update Skill", func(t *testing.T) {
		updatedName := "Golang"
		updateDTO := ports.UpdateSkillInput{
			Name:     &updatedName,
			Category: StrPtr("Backend"),
		}
		body, _ := json.Marshal(updateDTO)
		url := fmt.Sprintf("%s/%d", constSkillBase, createdSkillID)
		req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedSkill ports.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &updatedSkill)
		assert.Equal(t, updatedName, updatedSkill.Name)
		assert.Equal(t, "Backend", updatedSkill.Category)
	})

	t.Run("Toggle Skill Status", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d/%s", constSkillBase, createdSkillID, "toggle-status")
		req, _ := http.NewRequest(http.MethodPatch, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var toggledSkill ports.SkillResponse
		json.Unmarshal(w.Body.Bytes(), &toggledSkill)
		assert.False(t, toggledSkill.Active) 
	})


	t.Run("Get Skills with Filters and Search", func(t *testing.T) {
		createDTO := ports.CreateSkillInput{
			Category: "Database",
			Name:     "PostgreSQL",
			Active:   BoolPtr(true), 
		}
		body, _ := json.Marshal(createDTO)
		req, _ := http.NewRequest(http.MethodPost, constSkillBase, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)


		reqActive, _ := http.NewRequest(http.MethodGet, constSkillBase+"?active=true", nil)
		reqActive.Header.Set("Authorization", "Bearer "+userToken)
		wActive := httptest.NewRecorder()
		router.ServeHTTP(wActive, reqActive)
		var activeSkills []ports.SkillResponse
		json.Unmarshal(wActive.Body.Bytes(), &activeSkills)

		assert.Equal(t, 1, len(activeSkills), "Test 1: Filter active=true should return only 1 skill")
		assert.Equal(t, "PostgreSQL", activeSkills[0].Name)

		reqCat, _ := http.NewRequest(http.MethodGet, constSkillBase+"?category=Backend", nil)
		reqCat.Header.Set("Authorization", "Bearer "+userToken)
		wCat := httptest.NewRecorder()
		router.ServeHTTP(wCat, reqCat)
		var catSkills []ports.SkillResponse
		json.Unmarshal(wCat.Body.Bytes(), &catSkills)

		assert.Equal(t, 1, len(catSkills), "Test 2: Filter category=Backend should return 1 skill")
		assert.Equal(t, "Golang", catSkills[0].Name)

		reqSearch, _ := http.NewRequest(http.MethodGet, constSkillBase+"?search=Post", nil)
		reqSearch.Header.Set("Authorization", "Bearer "+userToken)
		wSearch := httptest.NewRecorder()
		router.ServeHTTP(wSearch, reqSearch)
		var searchSkills []ports.SkillResponse
		json.Unmarshal(wSearch.Body.Bytes(), &searchSkills)

		assert.Equal(t, 1, len(searchSkills), "Test 3: Search=Post should return 1 skill")
		assert.Equal(t, "PostgreSQL", searchSkills[0].Name)
	})

	t.Run("Delete Skill - In Use (Fail)", func(t *testing.T) {
		testSkillInUse := models.UserSkill{
			SkillID:          createdSkillID,
			UserID:           authorUser.ID,
			ProficiencyLevel: models.ProficiencyBeginner,
		}
		result := testutil.TestDB.Create(&testSkillInUse)
		assert.NoError(t, result.Error)

		url := fmt.Sprintf("%s/%d", constSkillBase, createdSkillID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code, "Delete should fail if skill is in use")
	})

	t.Run("Delete Skill - Success", func(t *testing.T) {
		result := testutil.TestDB.Delete(&models.UserSkill{}, "skill_id = ?", createdSkillID)
		assert.NoError(t, result.Error)

		url := fmt.Sprintf("%s/%d", constSkillBase, createdSkillID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Verify Deletion", func(t *testing.T) {
		url := fmt.Sprintf("%s/%d", constSkillBase, createdSkillID)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
