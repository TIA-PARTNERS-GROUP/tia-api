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

func TestProjectRegionService_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	prService := services.NewProjectRegionService(testutil.TestDB)
	regionService := services.NewRegionService(testutil.TestDB)

	user := models.User{FirstName: "RegionManager", LoginEmail: "region@manager.com", Active: true}
	testutil.TestDB.Create(&user)
	project := models.Project{Name: "Regional Project", ManagedByUserID: user.ID}
	testutil.TestDB.Create(&project)
	regionService.SeedRegions(context.Background())

	addDTO_QLD := ports.AddProjectRegionInput{ProjectID: project.ID, RegionID: "QLD"}
	_, err := prService.AddRegionToProject(context.Background(), addDTO_QLD)
	assert.NoError(t, err)

	_, err = prService.AddRegionToProject(context.Background(), addDTO_QLD)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrRegionAlreadyAdded, err)

	addDTO_NSW := ports.AddProjectRegionInput{ProjectID: project.ID, RegionID: "NSW"}
	_, err = prService.AddRegionToProject(context.Background(), addDTO_NSW)
	assert.NoError(t, err)

	regionsForProject, err := prService.GetRegionsForProject(context.Background(), project.ID)
	assert.NoError(t, err)
	assert.Len(t, regionsForProject, 2)

	assert.Equal(t, "NSW", regionsForProject[0].RegionID)
	assert.Equal(t, "New South Wales", regionsForProject[0].Region.Name)
	assert.Equal(t, "QLD", regionsForProject[1].RegionID)
	assert.Equal(t, "Queensland", regionsForProject[1].Region.Name)

	err = prService.RemoveRegionFromProject(context.Background(), project.ID, "QLD")
	assert.NoError(t, err)

	regionsForProject, _ = prService.GetRegionsForProject(context.Background(), project.ID)
	assert.Len(t, regionsForProject, 1)
	assert.Equal(t, "NSW", regionsForProject[0].RegionID)

	err = prService.RemoveRegionFromProject(context.Background(), project.ID, "QLD")
	assert.Error(t, err)
	assert.Equal(t, ports.ErrProjectRegionNotFound, err)
}
