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

func TestProjectApplicantService_Integration_ApplyAndWithdraw(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	applicantService := services.NewProjectApplicantService(testutil.TestDB)

	manager := models.User{FirstName: "Manager", LoginEmail: "manager@apply.com", Active: true}
	testutil.TestDB.Create(&manager)
	applicant := models.User{FirstName: "Applicant", LoginEmail: "applicant@apply.com", Active: true}
	testutil.TestDB.Create(&applicant)
	project := models.Project{Name: "Apply Project", ManagedByUserID: manager.ID}
	testutil.TestDB.Create(&project)

	applyDTO := ports.ApplyToProjectInput{
		ProjectID: project.ID,
		UserID:    applicant.ID,
	}
	application, err := applicantService.ApplyToProject(context.Background(), applyDTO)
	assert.NoError(t, err)
	assert.NotNil(t, application)

	_, err = applicantService.ApplyToProject(context.Background(), applyDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrAlreadyApplied, err)

	err = applicantService.WithdrawApplication(context.Background(), project.ID, applicant.ID)
	assert.NoError(t, err)

	err = applicantService.WithdrawApplication(context.Background(), project.ID, applicant.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrApplicationNotFound, err)
}

func TestProjectApplicantService_Integration_Getters(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	applicantService := services.NewProjectApplicantService(testutil.TestDB)

	manager := models.User{FirstName: "Manager", LoginEmail: "manager@get.com", Active: true}
	testutil.TestDB.Create(&manager)
	applicant1 := models.User{FirstName: "Applicant1", LoginEmail: "app1@get.com", Active: true}
	testutil.TestDB.Create(&applicant1)
	applicant2 := models.User{FirstName: "Applicant2", LoginEmail: "app2@get.com", Active: true}
	testutil.TestDB.Create(&applicant2)

	project1 := models.Project{Name: "Project 1", ManagedByUserID: manager.ID}
	testutil.TestDB.Create(&project1)
	project2 := models.Project{Name: "Project 2", ManagedByUserID: manager.ID}
	testutil.TestDB.Create(&project2)

	testutil.TestDB.Create(&models.ProjectApplicant{ProjectID: project1.ID, UserID: applicant1.ID})
	testutil.TestDB.Create(&models.ProjectApplicant{ProjectID: project1.ID, UserID: applicant2.ID})
	testutil.TestDB.Create(&models.ProjectApplicant{ProjectID: project2.ID, UserID: applicant1.ID})

	t.Run("Get Applicants For Project", func(t *testing.T) {
		applicants, err := applicantService.GetApplicantsForProject(context.Background(), project1.ID)
		assert.NoError(t, err)
		assert.Len(t, applicants, 2)
		assert.Equal(t, "Applicant1", applicants[0].User.FirstName)
	})

	t.Run("Get Applications For User", func(t *testing.T) {
		applications, err := applicantService.GetApplicationsForUser(context.Background(), applicant1.ID)
		assert.NoError(t, err)
		assert.Len(t, applications, 2)
		assert.Equal(t, "Project 1", applications[0].Project.Name)
	})
}
