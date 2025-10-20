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

func TestProjectService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectService := services.NewProjectService(testutil.TestDB)

	manager := models.User{FirstName: "Manager", LoginEmail: "manager@proj.com", Active: true}
	testutil.TestDB.Create(&manager)
	biz := models.Business{Name: "ProjectBiz", OperatorUserID: manager.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Create(&biz)

	createDTO := ports.CreateProjectInput{
		ManagedByUserID: manager.ID,
		BusinessID:      &biz.ID,
		Name:            "New Cloud Project",
		ProjectStatus:   models.ProjectStatusPlanning,
	}
	createdProject, err := projectService.CreateProject(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdProject)
	assert.Equal(t, "New Cloud Project", createdProject.Name)
	assert.Equal(t, manager.ID, createdProject.ManagedByUserID)
	assert.Len(t, createdProject.ProjectMembers, 1)
	assert.Equal(t, manager.ID, createdProject.ProjectMembers[0].UserID)
	assert.Equal(t, models.ProjectMemberRoleManager, createdProject.ProjectMembers[0].Role)

	fetchedProject, err := projectService.GetProjectByID(context.Background(), createdProject.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedProject)
	assert.Equal(t, "New Cloud Project", fetchedProject.Name)

	assert.Equal(t, manager.ID, fetchedProject.ManagingUser.ID)
	assert.NotNil(t, fetchedProject.Business)
	assert.Equal(t, biz.ID, fetchedProject.Business.ID)
}

func TestProjectService_Integration_MemberManagement(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	projectService := services.NewProjectService(testutil.TestDB)

	manager := models.User{FirstName: "Manager", LoginEmail: "manager@proj.com", Active: true}
	testutil.TestDB.Create(&manager)
	contributor := models.User{FirstName: "Contributor", LoginEmail: "contrib@proj.com", Active: true}
	testutil.TestDB.Create(&contributor)

	project := models.Project{Name: "Member Project", ManagedByUserID: manager.ID}
	testutil.TestDB.Create(&project)

	addDTO := ports.AddMemberInput{
		UserID: contributor.ID,
		Role:   models.ProjectMemberRoleContributor,
	}
	addedMember, err := projectService.AddMember(context.Background(), project.ID, addDTO)
	assert.NoError(t, err)
	assert.NotNil(t, addedMember)
	assert.Equal(t, contributor.ID, addedMember.UserID)
	assert.Equal(t, "Contributor", addedMember.User.FirstName)

	_, err = projectService.AddMember(context.Background(), project.ID, addDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrMemberAlreadyExists, err)

	err = projectService.RemoveMember(context.Background(), project.ID, contributor.ID)
	assert.NoError(t, err)

	err = projectService.RemoveMember(context.Background(), project.ID, contributor.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrMemberNotFound, err)
}
