package main
import (
	"context"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestSkillService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	skillService := services.NewSkillService(testutil.TestDB)
	createDTO := ports.CreateSkillInput{
		Category: "Programming",
		Name:     "Golang",
	}
	createdSkill, err := skillService.CreateSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdSkill)
	assert.Equal(t, "Golang", createdSkill.Name)
	assert.True(t, createdSkill.Active)
	fetchedSkill, err := skillService.GetSkillByID(context.Background(), createdSkill.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedSkill)
	assert.Equal(t, "Golang", fetchedSkill.Name)
}
func TestSkillService_Integration_CreateDuplicate(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	skillService := services.NewSkillService(testutil.TestDB)
	testutil.TestDB.Create(&models.Skill{Name: "Duplicate Skill", Category: "Test"})
	createDTO := ports.CreateSkillInput{Name: "Duplicate Skill", Category: "Test"}
	_, err := skillService.CreateSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrSkillNameExists, err)
}
func TestSkillService_Integration_DeleteSkill(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	skillService := services.NewSkillService(testutil.TestDB)
	t.Run("Success - No Dependencies", func(t *testing.T) {
		skill := models.Skill{Name: "Deletable", Category: "Test"}
		testutil.TestDB.Create(&skill)
		err := skillService.DeleteSkill(context.Background(), skill.ID)
		assert.NoError(t, err)
	})
	t.Run("Failure - Skill In Use", func(t *testing.T) {
		user := models.User{FirstName: "SkillUser", LoginEmail: "skill@user.com"}
		testutil.TestDB.Create(&user)
		skill := models.Skill{Name: "InUse Skill", Category: "Test"}
		testutil.TestDB.Create(&skill)
		testutil.TestDB.Create(&models.UserSkill{UserID: user.ID, SkillID: skill.ID})
		err := skillService.DeleteSkill(context.Background(), skill.ID)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrSkillInUse, err)
	})
}
func TestSkillService_Integration_GetSkillsFiltered(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	skillService := services.NewSkillService(testutil.TestDB)
	testutil.TestDB.Model(&models.Skill{}).Create(map[string]interface{}{"Name": "Golang", "Category": "Backend", "Active": true})
	testutil.TestDB.Model(&models.Skill{}).Create(map[string]interface{}{"Name": "TypeScript", "Category": "Frontend", "Active": true})
	testutil.TestDB.Model(&models.Skill{}).Create(map[string]interface{}{"Name": "Java", "Category": "Backend", "Active": false})
	t.Run("Filter by Category", func(t *testing.T) {
		category := "Backend"
		filters := ports.SkillsFilter{Category: &category}
		skills, err := skillService.GetSkills(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, skills, 2)
	})
	t.Run("Filter by Active status", func(t *testing.T) {
		active := false
		filters := ports.SkillsFilter{Active: &active}
		skills, err := skillService.GetSkills(context.Background(), filters)
		assert.NoError(t, err)
		if assert.Len(t, skills, 1) {
			assert.Equal(t, "Java", skills[0].Name)
		}
	})
	t.Run("Filter by Search term", func(t *testing.T) {
		search := "Script"
		filters := ports.SkillsFilter{Search: &search}
		skills, err := skillService.GetSkills(context.Background(), filters)
		assert.NoError(t, err)
		if assert.Len(t, skills, 1) {
			assert.Equal(t, "TypeScript", skills[0].Name)
		}
	})
}
