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

func TestUserSubscriptionService_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	userSubService := services.NewUserSubscriptionService(testutil.TestDB)

	user := models.User{FirstName: "Subscriber", LoginEmail: "sub@user.com", Active: true}
	testutil.TestDB.Create(&user)

	days := 30
	plan := models.Subscription{Name: "Monthly Plan", Price: 9.99, ValidDays: &days}
	testutil.TestDB.Create(&plan)

	createDTO := ports.CreateUserSubscriptionInput{
		UserID:         user.ID,
		SubscriptionID: plan.ID,
		IsTrial:        true,
	}
	createdSub, err := userSubService.CreateUserSubscription(context.Background(), createDTO)

	assert.NoError(t, err)
	assert.NotNil(t, createdSub)
	assert.Equal(t, user.ID, createdSub.UserID)
	assert.Equal(t, plan.ID, createdSub.SubscriptionID)
	assert.True(t, createdSub.IsTrial)

	expectedExpiry := time.Now().AddDate(0, 0, 30)
	assert.WithinDuration(t, expectedExpiry, createdSub.DateTo, 5*time.Second)

	fetchedSub, err := userSubService.GetUserSubscriptionByID(context.Background(), createdSub.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedSub)
	assert.Equal(t, "Subscriber", fetchedSub.User.FirstName)
	assert.Equal(t, "Monthly Plan", fetchedSub.Subscription.Name)

	userSubs, err := userSubService.GetSubscriptionsForUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, userSubs, 1)

	err = userSubService.CancelSubscription(context.Background(), createdSub.ID)
	assert.NoError(t, err)

	userSubs, err = userSubService.GetSubscriptionsForUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, userSubs, 0)
}
