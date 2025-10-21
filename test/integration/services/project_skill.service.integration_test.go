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
func TestProjectSkillService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusPlanning,
	}
	testutil.TestDB.Create(&project)
	skill := models.Skill{Name: "Go Programming", Category: "Programming", Active: true}
	testutil.TestDB.Create(&skill)
	createDTO := ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    skill.ID,
		Importance: models.SkillImportanceRequired,
	}
	createdProjectSkill, err := projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdProjectSkill)
	assert.Equal(t, project.ID, createdProjectSkill.ProjectID)
	assert.Equal(t, skill.ID, createdProjectSkill.SkillID)
	assert.Equal(t, models.SkillImportanceRequired, createdProjectSkill.Importance)
	assert.NotNil(t, createdProjectSkill.Skill)
	assert.Equal(t, skill.Name, createdProjectSkill.Skill.Name)
	assert.NotNil(t, createdProjectSkill.Project)
	assert.Equal(t, project.Name, createdProjectSkill.Project.Name)
	fetchedProjectSkill, err := projectSkillService.GetProjectSkill(context.Background(), project.ID, skill.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedProjectSkill)
	assert.Equal(t, createdProjectSkill.Importance, fetchedProjectSkill.Importance)
}
func TestProjectSkillService_Integration_DuplicatePrevention(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager2@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 2",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	skill := models.Skill{Name: "Python Programming", Category: "Programming", Active: true}
	testutil.TestDB.Create(&skill)
	createDTO := ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    skill.ID,
		Importance: models.SkillImportancePreferred,
	}
	_, err := projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectSkillAlreadyExists, err)
}
func TestProjectSkillService_Integration_UpdateAndDelete(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager3@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 3",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	skill := models.Skill{Name: "JavaScript", Category: "Programming", Active: true}
	testutil.TestDB.Create(&skill)
	createDTO := ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    skill.ID,
		Importance: models.SkillImportanceOptional,
	}
	_, err := projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	requiredImportance := models.SkillImportanceRequired
	updateDTO := ports.UpdateProjectSkillInput{
		Importance: &requiredImportance,
	}
	updatedProjectSkill, err := projectSkillService.UpdateProjectSkill(context.Background(), project.ID, skill.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedProjectSkill)
	assert.Equal(t, models.SkillImportanceRequired, updatedProjectSkill.Importance)
	projectSkills, err := projectSkillService.GetProjectSkills(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Len(t, projectSkills, 1)
	assert.Equal(t, skill.ID, projectSkills[0].SkillID)
	assert.Equal(t, models.SkillImportanceRequired, projectSkills[0].Importance)
	err = projectSkillService.RemoveProjectSkill(context.Background(), project.ID, skill.ID)
	assert.NoError(t, err)
	_, err = projectSkillService.GetProjectSkill(context.Background(), project.ID, skill.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectSkillNotFound, err)
	projectSkills, err = projectSkillService.GetProjectSkills(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Len(t, projectSkills, 0)
}
func TestProjectSkillService_Integration_GetProjectsBySkill(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	skill1 := models.Skill{Name: "React", Category: "Frontend", Active: true}
	testutil.TestDB.Create(&skill1)
	skill2 := models.Skill{Name: "Node.js", Category: "Backend", Active: true}
	testutil.TestDB.Create(&skill2)
	manager1 := models.User{FirstName: "Manager1", LoginEmail: "manager1@test.com", Active: true}
	testutil.TestDB.Create(&manager1)
	manager2 := models.User{FirstName: "Manager2", LoginEmail: "manager2@test.com", Active: true}
	testutil.TestDB.Create(&manager2)
	project1 := models.Project{
		Name:            "Frontend Project",
		ManagedByUserID: manager1.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project1)
	project2 := models.Project{
		Name:            "Backend Project",
		ManagedByUserID: manager2.ID,
		ProjectStatus:   models.ProjectStatusPlanning,
	}
	testutil.TestDB.Create(&project2)
	project3 := models.Project{
		Name:            "Fullstack Project",
		ManagedByUserID: manager1.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project3)
	_, err := projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project1.ID,
		SkillID:    skill1.ID,
		Importance: models.SkillImportanceRequired,
	})
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project2.ID,
		SkillID:    skill2.ID,
		Importance: models.SkillImportanceRequired,
	})
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project3.ID,
		SkillID:    skill1.ID,
		Importance: models.SkillImportancePreferred,
	})
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project3.ID,
		SkillID:    skill2.ID,
		Importance: models.SkillImportanceRequired,
	})
	assert.NoError(t, err)
	reactProjects, err := projectSkillService.GetProjectsBySkill(context.Background(), skill1.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, reactProjects, 2)
	nodeProjects, err := projectSkillService.GetProjectsBySkill(context.Background(), skill2.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, nodeProjects, 2)
	requiredImportance := models.SkillImportanceRequired
	requiredReactProjects, err := projectSkillService.GetProjectsBySkill(context.Background(), skill1.ID, &requiredImportance)
	assert.NoError(t, err)
	assert.Len(t, requiredReactProjects, 1)
	assert.Equal(t, project1.ID, requiredReactProjects[0].ProjectID)
}
func TestProjectSkillService_Integration_Validation(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	createDTO := ports.CreateProjectSkillInput{
		ProjectID:  999,
		SkillID:    1,
		Importance: models.SkillImportanceRequired,
	}
	_, err := projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectNotFound, err)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager4@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 4",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	createDTO = ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    999,
		Importance: models.SkillImportanceRequired,
	}
	_, err = projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrSkillNotFound, err)
}
func TestProjectSkillService_Integration_GetSkillsByImportance(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager5@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Complex Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	requiredSkill := models.Skill{Name: "Required Skill", Category: "Core", Active: true}
	testutil.TestDB.Create(&requiredSkill)
	preferredSkill := models.Skill{Name: "Preferred Skill", Category: "Nice to have", Active: true}
	testutil.TestDB.Create(&preferredSkill)
	optionalSkill := models.Skill{Name: "Optional Skill", Category: "Extra", Active: true}
	testutil.TestDB.Create(&optionalSkill)
	_, err := projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    requiredSkill.ID,
		Importance: models.SkillImportanceRequired,
	})
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    preferredSkill.ID,
		Importance: models.SkillImportancePreferred,
	})
	assert.NoError(t, err)
	_, err = projectSkillService.AddProjectSkill(context.Background(), ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    optionalSkill.ID,
		Importance: models.SkillImportanceOptional,
	})
	assert.NoError(t, err)
	requiredSkills, err := projectSkillService.GetSkillsByImportance(context.Background(), project.ID, models.SkillImportanceRequired)
	assert.NoError(t, err)
	assert.Len(t, requiredSkills, 1)
	assert.Equal(t, requiredSkill.ID, requiredSkills[0].SkillID)
	preferredSkills, err := projectSkillService.GetSkillsByImportance(context.Background(), project.ID, models.SkillImportancePreferred)
	assert.NoError(t, err)
	assert.Len(t, preferredSkills, 1)
	assert.Equal(t, preferredSkill.ID, preferredSkills[0].SkillID)
	optionalSkills, err := projectSkillService.GetSkillsByImportance(context.Background(), project.ID, models.SkillImportanceOptional)
	assert.NoError(t, err)
	assert.Len(t, optionalSkills, 1)
	assert.Equal(t, optionalSkill.ID, optionalSkills[0].SkillID)
}
func TestProjectSkillService_Integration_NonExistentProjectSkill(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	_, err := projectSkillService.GetProjectSkill(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectSkillNotFound, err)
	requiredImportance := models.SkillImportanceRequired
	updateDTO := ports.UpdateProjectSkillInput{
		Importance: &requiredImportance,
	}
	_, err = projectSkillService.UpdateProjectSkill(context.Background(), 999, 999, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectSkillNotFound, err)
	err = projectSkillService.RemoveProjectSkill(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectSkillNotFound, err)
}
func TestProjectSkillService_Integration_UpdateNoData(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectSkillService := services.NewProjectSkillService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager6@projectskill.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 6",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	skill := models.Skill{Name: "Vue.js", Category: "Frontend", Active: true}
	testutil.TestDB.Create(&skill)
	createDTO := ports.CreateProjectSkillInput{
		ProjectID:  project.ID,
		SkillID:    skill.ID,
		Importance: models.SkillImportanceRequired,
	}
	_, err := projectSkillService.AddProjectSkill(context.Background(), createDTO)
	assert.NoError(t, err)
	updateDTO := ports.UpdateProjectSkillInput{}
	_, err = projectSkillService.UpdateProjectSkill(context.Background(), project.ID, skill.ID, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrNoUpdateData, err)
}
