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
func TestProjectMemberService_Integration_AddAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusPlanning,
	}
	testutil.TestDB.Create(&project)
	user := models.User{FirstName: "Member", LoginEmail: "member@project.com", Active: true}
	testutil.TestDB.Create(&user)
	addDTO := ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    user.ID,
		Role:      models.ProjectMemberRoleContributor,
	}
	createdMember, err := projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdMember)
	assert.Equal(t, project.ID, createdMember.ProjectID)
	assert.Equal(t, user.ID, createdMember.UserID)
	assert.Equal(t, models.ProjectMemberRoleContributor, createdMember.Role)
	assert.NotNil(t, createdMember.Project)
	assert.Equal(t, project.Name, createdMember.Project.Name)
	assert.NotNil(t, createdMember.User)
	assert.Equal(t, user.FirstName, createdMember.User.FirstName)
	fetchedMember, err := projectMemberService.GetProjectMember(context.Background(), project.ID, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedMember)
	assert.Equal(t, createdMember.Role, fetchedMember.Role)
	assert.Equal(t, createdMember.JoinedAt, fetchedMember.JoinedAt)
}
func TestProjectMemberService_Integration_DuplicatePrevention(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager2@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 2",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	user := models.User{FirstName: "Member", LoginEmail: "member2@project.com", Active: true}
	testutil.TestDB.Create(&user)
	addDTO := ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    user.ID,
		Role:      models.ProjectMemberRoleReviewer,
	}
	_, err := projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectMemberAlreadyExists, err)
}
func TestProjectMemberService_Integration_UpdateAndRemove(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager3@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 3",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	user := models.User{FirstName: "Member", LoginEmail: "member3@project.com", Active: true}
	testutil.TestDB.Create(&user)
	addDTO := ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    user.ID,
		Role:      models.ProjectMemberRoleContributor,
	}
	_, err := projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.NoError(t, err)
	reviewerRole := models.ProjectMemberRoleReviewer
	updateDTO := ports.UpdateProjectMemberRoleInput{
		Role: reviewerRole,
	}
	updatedMember, err := projectMemberService.UpdateProjectMemberRole(context.Background(), project.ID, user.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedMember)
	assert.Equal(t, models.ProjectMemberRoleReviewer, updatedMember.Role)
	members, err := projectMemberService.GetProjectMembers(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Len(t, members, 1)
	assert.Equal(t, user.ID, members[0].UserID)
	assert.Equal(t, models.ProjectMemberRoleReviewer, members[0].Role)
	err = projectMemberService.RemoveProjectMember(context.Background(), project.ID, user.ID)
	assert.NoError(t, err)
	_, err = projectMemberService.GetProjectMember(context.Background(), project.ID, user.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectMemberNotFound, err)
	members, err = projectMemberService.GetProjectMembers(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Len(t, members, 0)
}
func TestProjectMemberService_Integration_GetProjectsByUser(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager1 := models.User{FirstName: "Manager1", LoginEmail: "manager1@test.com", Active: true}
	testutil.TestDB.Create(&manager1)
	manager2 := models.User{FirstName: "Manager2", LoginEmail: "manager2@test.com", Active: true}
	testutil.TestDB.Create(&manager2)
	project1 := models.Project{
		Name:            "Project 1",
		ManagedByUserID: manager1.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project1)
	project2 := models.Project{
		Name:            "Project 2",
		ManagedByUserID: manager2.ID,
		ProjectStatus:   models.ProjectStatusPlanning,
	}
	testutil.TestDB.Create(&project2)
	project3 := models.Project{
		Name:            "Project 3",
		ManagedByUserID: manager1.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project3)
	member := models.User{FirstName: "MultiProject", LoginEmail: "multiproject@test.com", Active: true}
	testutil.TestDB.Create(&member)
	_, err := projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project1.ID,
		UserID:    member.ID,
		Role:      models.ProjectMemberRoleContributor,
	})
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project2.ID,
		UserID:    member.ID,
		Role:      models.ProjectMemberRoleReviewer,
	})
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project3.ID,
		UserID:    member.ID,
		Role:      models.ProjectMemberRoleManager,
	})
	assert.NoError(t, err)
	userProjects, err := projectMemberService.GetProjectsByUser(context.Background(), member.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, userProjects, 3)
	managerRole := models.ProjectMemberRoleManager
	managerProjects, err := projectMemberService.GetProjectsByUser(context.Background(), member.ID, &managerRole)
	assert.NoError(t, err)
	assert.Len(t, managerProjects, 1)
	assert.Equal(t, project3.ID, managerProjects[0].ProjectID)
}
func TestProjectMemberService_Integration_Validation(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	addDTO := ports.AddProjectMemberInput{
		ProjectID: 999,
		UserID:    1,
		Role:      models.ProjectMemberRoleContributor,
	}
	_, err := projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectNotFound, err)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager4@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Test Project 4",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	addDTO = ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    999,
		Role:      models.ProjectMemberRoleContributor,
	}
	_, err = projectMemberService.AddProjectMember(context.Background(), addDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserNotFound, err)
}
func TestProjectMemberService_Integration_GetMembersByRole(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager5@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Team Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	contributor1 := models.User{FirstName: "Contributor1", LoginEmail: "contributor1@test.com", Active: true}
	testutil.TestDB.Create(&contributor1)
	contributor2 := models.User{FirstName: "Contributor2", LoginEmail: "contributor2@test.com", Active: true}
	testutil.TestDB.Create(&contributor2)
	reviewer := models.User{FirstName: "Reviewer", LoginEmail: "reviewer@test.com", Active: true}
	testutil.TestDB.Create(&reviewer)
	manager2 := models.User{FirstName: "Manager2", LoginEmail: "manager2@test.com", Active: true}
	testutil.TestDB.Create(&manager2)
	_, err := projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    contributor1.ID,
		Role:      models.ProjectMemberRoleContributor,
	})
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    contributor2.ID,
		Role:      models.ProjectMemberRoleContributor,
	})
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    reviewer.ID,
		Role:      models.ProjectMemberRoleReviewer,
	})
	assert.NoError(t, err)
	_, err = projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    manager2.ID,
		Role:      models.ProjectMemberRoleManager,
	})
	assert.NoError(t, err)
	contributors, err := projectMemberService.GetMembersByRole(context.Background(), project.ID, models.ProjectMemberRoleContributor)
	assert.NoError(t, err)
	assert.Len(t, contributors, 2)
	reviewers, err := projectMemberService.GetMembersByRole(context.Background(), project.ID, models.ProjectMemberRoleReviewer)
	assert.NoError(t, err)
	assert.Len(t, reviewers, 1)
	assert.Equal(t, reviewer.ID, reviewers[0].UserID)
	managers, err := projectMemberService.GetMembersByRole(context.Background(), project.ID, models.ProjectMemberRoleManager)
	assert.NoError(t, err)
	assert.Len(t, managers, 1)
	assert.Equal(t, manager2.ID, managers[0].UserID)
}
func TestProjectMemberService_Integration_CannotRemoveManager(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager6@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Manager Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	_, err := projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    manager.ID,
		Role:      models.ProjectMemberRoleManager,
	})
	assert.NoError(t, err)
	err = projectMemberService.RemoveProjectMember(context.Background(), project.ID, manager.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrCannotRemoveManager, err)
}
func TestProjectMemberService_Integration_IsUserProjectMember(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	manager := models.User{FirstName: "Manager", LoginEmail: "manager7@projectmember.com", Active: true}
	testutil.TestDB.Create(&manager)
	project := models.Project{
		Name:            "Membership Test Project",
		ManagedByUserID: manager.ID,
		ProjectStatus:   models.ProjectStatusActive,
	}
	testutil.TestDB.Create(&project)
	member := models.User{FirstName: "Member", LoginEmail: "member7@test.com", Active: true}
	testutil.TestDB.Create(&member)
	nonMember := models.User{FirstName: "NonMember", LoginEmail: "nonmember@test.com", Active: true}
	testutil.TestDB.Create(&nonMember)
	_, err := projectMemberService.AddProjectMember(context.Background(), ports.AddProjectMemberInput{
		ProjectID: project.ID,
		UserID:    member.ID,
		Role:      models.ProjectMemberRoleContributor,
	})
	assert.NoError(t, err)
	isMember, err := projectMemberService.IsUserProjectMember(context.Background(), project.ID, member.ID)
	assert.NoError(t, err)
	assert.True(t, isMember)
	isMember, err = projectMemberService.IsUserProjectMember(context.Background(), project.ID, nonMember.ID)
	assert.NoError(t, err)
	assert.False(t, isMember)
}
func TestProjectMemberService_Integration_NonExistentProjectMember(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectMemberService := services.NewProjectMemberService(testutil.TestDB)
	_, err := projectMemberService.GetProjectMember(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectMemberNotFound, err)
	contributorRole := models.ProjectMemberRoleContributor
	updateDTO := ports.UpdateProjectMemberRoleInput{
		Role: contributorRole,
	}
	_, err = projectMemberService.UpdateProjectMemberRole(context.Background(), 999, 999, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectMemberNotFound, err)
	err = projectMemberService.RemoveProjectMember(context.Background(), 999, 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectMemberNotFound, err)
}
