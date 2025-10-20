package main

import (
	"context"
	"testing"
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionService_Integration_PlanManagement(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	subService := services.NewSubscriptionService(testutil.TestDB)

	days := 30
	createDTO := ports.CreateSubscriptionInput{
		Name:      "Monthly Plan",
		Price:     9.99,
		ValidDays: &days,
	}
	createdPlan, err := subService.CreateSubscription(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdPlan)
	assert.Equal(t, "Monthly Plan", createdPlan.Name)
	assert.Equal(t, 30, *createdPlan.ValidDays)
}

func TestSubscriptionService_Integration_UserSubscription(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	subService := services.NewSubscriptionService(testutil.TestDB)

	user := models.User{FirstName: "Sub", LoginEmail: "sub@user.com", Active: true}
	testutil.TestDB.Create(&user)

	months := 1
	plan := models.Subscription{Name: "Yearly Plan", Price: 99.99, ValidMonths: &months}
	testutil.TestDB.Create(&plan)

	subscribeDTO := ports.UserSubscribeInput{
		UserID:         user.ID,
		SubscriptionID: plan.ID,
	}
	userSub, err := subService.SubscribeUser(context.Background(), subscribeDTO)
	assert.NoError(t, err)
	assert.NotNil(t, userSub)
	assert.Equal(t, user.ID, userSub.UserID)
	assert.Equal(t, plan.ID, userSub.SubscriptionID)

	expectedExpiry := time.Now().AddDate(0, 1, 0)
	assert.WithinDuration(t, expectedExpiry, userSub.DateTo, 5*time.Second)

	fetchedSub, err := subService.GetUserSubscription(context.Background(), userSub.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedSub)
	assert.Equal(t, user.ID, fetchedSub.User.ID)
	assert.Equal(t, plan.Name, fetchedSub.Subscription.Name)
}
