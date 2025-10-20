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

func TestBusinessService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)

	businessService := services.NewBusinessService(testutil.TestDB)
	operator := models.User{FirstName: "Biz", LoginEmail: "biz@owner.com", Active: true}

	testutil.TestDB.Create(&operator)

	createDTO := ports.CreateBusinessInput{
		OperatorUserID:   operator.ID,
		Name:             "My New Business",
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
	}
	createdBusiness, err := businessService.CreateBusiness(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdBusiness)
	assert.Equal(t, "My New Business", createdBusiness.Name)
	assert.Equal(t, operator.ID, createdBusiness.OperatorUserID)

	fetchedBusiness, err := businessService.GetBusinessByID(context.Background(), createdBusiness.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedBusiness)
	assert.Equal(t, "My New Business", fetchedBusiness.Name)

	assert.NotZero(t, fetchedBusiness.OperatorUser.ID)
	assert.Equal(t, "Biz", fetchedBusiness.OperatorUser.FirstName)
}

func TestBusinessService_Integration_DeleteBusiness(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessService := services.NewBusinessService(testutil.TestDB)

	operator := models.User{FirstName: "BizDel", LoginEmail: "bizdel@owner.com", Active: true}
	testutil.TestDB.Create(&operator)

	t.Run("Success - No Dependencies", func(t *testing.T) {
		businessToDelete := models.Business{Name: "Deletable", OperatorUserID: operator.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
		testutil.TestDB.Create(&businessToDelete)

		err := businessService.DeleteBusiness(context.Background(), businessToDelete.ID)
		assert.NoError(t, err)
	})

	t.Run("Failure - Business In Use", func(t *testing.T) {
		businessInUse := models.Business{Name: "In Use", OperatorUserID: operator.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
		testutil.TestDB.Create(&businessInUse)
		project := models.Project{Name: "Test Project", ManagedByUserID: operator.ID, BusinessID: &businessInUse.ID}
		testutil.TestDB.Create(&project)

		err := businessService.DeleteBusiness(context.Background(), businessInUse.ID)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrBusinessInUse, err)
	})
}

func TestBusinessService_Integration_GetUserBusinesses(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessService := services.NewBusinessService(testutil.TestDB)

	user1 := models.User{FirstName: "User1", LoginEmail: "user1@biz.com", Active: true}
	testutil.TestDB.Create(&user1)
	user2 := models.User{FirstName: "User2", LoginEmail: "user2@biz.com", Active: true}
	testutil.TestDB.Create(&user2)

	testutil.TestDB.Create(&models.Business{Name: "Biz 1", OperatorUserID: user1.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"})
	testutil.TestDB.Create(&models.Business{Name: "Biz 2", OperatorUserID: user1.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"})
	testutil.TestDB.Create(&models.Business{Name: "Biz 3", OperatorUserID: user2.ID, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"})

	user1Businesses, err := businessService.GetUserBusinesses(context.Background(), user1.ID)
	assert.NoError(t, err)
	assert.Len(t, user1Businesses, 2)
}

func TestBusinessService_Integration_ToggleBusinessStatus(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessService := services.NewBusinessService(testutil.TestDB)

	operator := models.User{FirstName: "ToggleUser", LoginEmail: "toggle@user.com", Active: true}
	testutil.TestDB.Create(&operator)
	business := models.Business{Name: "Toggle Biz", OperatorUserID: operator.ID, Active: true, BusinessType: "Other", BusinessCategory: "Mixed", BusinessPhase: "Growth"}
	testutil.TestDB.Select("*").Create(&business) // Ensure 'Active' is set

	updatedBiz, err := businessService.ToggleBusinessStatus(context.Background(), business.ID)
	assert.NoError(t, err)
	assert.False(t, updatedBiz.Active)

	updatedBiz, err = businessService.ToggleBusinessStatus(context.Background(), business.ID)
	assert.NoError(t, err)
	assert.True(t, updatedBiz.Active)
}
