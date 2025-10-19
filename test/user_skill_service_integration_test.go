package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/stretchr/testify/assert"
)

func TestUserSkillService_Integration_CreateAndGet(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	user := models.User{FirstName: "Test", LoginEmail: "test@userskill.com", Active: true}
	testDB.Create(&user)

	skill := models.Skill{Name: "Go Programming", Category: "Programming", Active: true}
	testDB.Create(&skill)

	createDTO := ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          skill.ID,
		ProficiencyLevel: models.ProficiencyAdvanced,
	}

	createdUserSkill, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdUserSkill)
	assert.Equal(t, user.ID, createdUserSkill.UserID)
	assert.Equal(t, skill.ID, createdUserSkill.SkillID)
	assert.Equal(t, models.ProficiencyAdvanced, createdUserSkill.ProficiencyLevel)
	assert.NotNil(t, createdUserSkill.Skill)
	assert.Equal(t, skill.Name, createdUserSkill.Skill.Name)
	assert.NotNil(t, createdUserSkill.User)
	assert.Equal(t, user.FirstName, createdUserSkill.User.FirstName)

	fetchedUserSkill, err := userSkillService.GetUserSkill(context.Background(), user.ID, skill.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedUserSkill)
	assert.Equal(t, createdUserSkill.ProficiencyLevel, fetchedUserSkill.ProficiencyLevel)
	assert.Equal(t, createdUserSkill.CreatedAt, fetchedUserSkill.CreatedAt)
}

func TestUserSkillService_Integration_DuplicatePrevention(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	user := models.User{FirstName: "Test", LoginEmail: "test2@userskill.com", Active: true}
	testDB.Create(&user)

	skill := models.Skill{Name: "Python Programming", Category: "Programming", Active: true}
	testDB.Create(&skill)

	createDTO := ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          skill.ID,
		ProficiencyLevel: models.ProficiencyIntermediate,
	}

	_, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.NoError(t, err)

	_, err = userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserSkillAlreadyExists, err)
}

func TestUserSkillService_Integration_UpdateAndDelete(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	user := models.User{FirstName: "Test", LoginEmail: "test3@userskill.com", Active: true}
	testDB.Create(&user)

	skill := models.Skill{Name: "JavaScript", Category: "Programming", Active: true}
	testDB.Create(&skill)

	createDTO := ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          skill.ID,
		ProficiencyLevel: models.ProficiencyBeginner,
	}

	_, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.NoError(t, err)

	expertProficiency := models.ProficiencyExpert
	updateDTO := ports.UpdateUserSkillInput{
		ProficiencyLevel: &expertProficiency,
	}

	updatedUserSkill, err := userSkillService.UpdateUserSkill(context.Background(), user.ID, skill.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUserSkill)
	assert.Equal(t, models.ProficiencyExpert, updatedUserSkill.ProficiencyLevel)

	userSkills, err := userSkillService.GetUserSkills(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, userSkills, 1)
	assert.Equal(t, skill.ID, userSkills[0].SkillID)
	assert.Equal(t, models.ProficiencyExpert, userSkills[0].ProficiencyLevel)

	err = userSkillService.RemoveUserSkill(context.Background(), user.ID, skill.ID)
	assert.NoError(t, err)

	_, err = userSkillService.GetUserSkill(context.Background(), user.ID, skill.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserSkillNotFound, err)

	userSkills, err = userSkillService.GetUserSkills(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, userSkills, 0)
}

func TestUserSkillService_Integration_GetUsersBySkill(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	skill1 := models.Skill{Name: "React", Category: "Frontend", Active: true}
	testDB.Create(&skill1)
	skill2 := models.Skill{Name: "Node.js", Category: "Backend", Active: true}
	testDB.Create(&skill2)

	user1 := models.User{FirstName: "Frontend", LoginEmail: "frontend@test.com", Active: true}
	testDB.Create(&user1)
	user2 := models.User{FirstName: "Backend", LoginEmail: "backend@test.com", Active: true}
	testDB.Create(&user2)
	user3 := models.User{FirstName: "Fullstack", LoginEmail: "fullstack@test.com", Active: true}
	testDB.Create(&user3)

	_, err := userSkillService.AddUserSkill(context.Background(), ports.CreateUserSkillInput{
		UserID:           user1.ID,
		SkillID:          skill1.ID,
		ProficiencyLevel: models.ProficiencyAdvanced,
	})
	assert.NoError(t, err)

	_, err = userSkillService.AddUserSkill(context.Background(), ports.CreateUserSkillInput{
		UserID:           user2.ID,
		SkillID:          skill2.ID,
		ProficiencyLevel: models.ProficiencyExpert,
	})
	assert.NoError(t, err)

	_, err = userSkillService.AddUserSkill(context.Background(), ports.CreateUserSkillInput{
		UserID:           user3.ID,
		SkillID:          skill1.ID,
		ProficiencyLevel: models.ProficiencyIntermediate,
	})
	assert.NoError(t, err)

	_, err = userSkillService.AddUserSkill(context.Background(), ports.CreateUserSkillInput{
		UserID:           user3.ID,
		SkillID:          skill2.ID,
		ProficiencyLevel: models.ProficiencyAdvanced,
	})
	assert.NoError(t, err)

	reactUsers, err := userSkillService.GetUsersBySkill(context.Background(), skill1.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, reactUsers, 2)

	nodeUsers, err := userSkillService.GetUsersBySkill(context.Background(), skill2.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, nodeUsers, 2)

	expertProficiency := models.ProficiencyExpert
	expertNodeUsers, err := userSkillService.GetUsersBySkill(context.Background(), skill2.ID, &expertProficiency)
	assert.NoError(t, err)
	assert.Len(t, expertNodeUsers, 1)
	assert.Equal(t, user2.ID, expertNodeUsers[0].UserID)
	assert.Equal(t, models.ProficiencyExpert, expertNodeUsers[0].ProficiencyLevel)

	intermediateProficiency := models.ProficiencyIntermediate
	intermediateReactUsers, err := userSkillService.GetUsersBySkill(context.Background(), skill1.ID, &intermediateProficiency)
	assert.NoError(t, err)
	assert.Len(t, intermediateReactUsers, 1)
	assert.Equal(t, user3.ID, intermediateReactUsers[0].UserID)
}

func TestUserSkillService_Integration_Validation(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	createDTO := ports.CreateUserSkillInput{
		UserID:           999,
		SkillID:          1,
		ProficiencyLevel: models.ProficiencyIntermediate,
	}

	_, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserNotFound, err)

	user := models.User{FirstName: "Test", LoginEmail: "test4@userskill.com", Active: true}
	testDB.Create(&user)

	createDTO = ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          999,
		ProficiencyLevel: models.ProficiencyIntermediate,
	}

	_, err = userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrSkillNotFound, err)
}

func TestUserSkillService_Integration_DefaultProficiency(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	user := models.User{FirstName: "Test", LoginEmail: "test5@userskill.com", Active: true}
	testDB.Create(&user)

	skill := models.Skill{Name: "TypeScript", Category: "Programming", Active: true}
	testDB.Create(&skill)

	createDTO := ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          skill.ID,
		ProficiencyLevel: models.ProficiencyIntermediate,
	}

	userSkill, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.Equal(t, models.ProficiencyIntermediate, userSkill.ProficiencyLevel)
}

func TestUserSkillService_Integration_UpdateNoData(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	user := models.User{FirstName: "Test", LoginEmail: "test6@userskill.com", Active: true}
	testDB.Create(&user)

	skill := models.Skill{Name: "Vue.js", Category: "Frontend", Active: true}
	testDB.Create(&skill)

	createDTO := ports.CreateUserSkillInput{
		UserID:           user.ID,
		SkillID:          skill.ID,
		ProficiencyLevel: models.ProficiencyAdvanced,
	}

	_, err := userSkillService.AddUserSkill(context.Background(), createDTO)
	assert.NoError(t, err)

	updateDTO := ports.UpdateUserSkillInput{}
	_, err = userSkillService.UpdateUserSkill(context.Background(), user.ID, skill.ID, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrNoUpdateData, err)
}

func TestUserSkillService_Integration_NonExistentUserSkill(t *testing.T) {
	cleanupTestDB(t)
	userSkillService := services.NewUserSkillService(testDB)

	_, err := userSkillService.GetUserSkill(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserSkillNotFound, err)

	expertProficiency := models.ProficiencyExpert
	updateDTO := ports.UpdateUserSkillInput{
		ProficiencyLevel: &expertProficiency,
	}
	_, err = userSkillService.UpdateUserSkill(context.Background(), 999, 999, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserSkillNotFound, err)

	err = userSkillService.RemoveUserSkill(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserSkillNotFound, err)
}
