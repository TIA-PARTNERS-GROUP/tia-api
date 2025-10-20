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

func TestDailyActivityEnrolmentService_Integration(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	enrolmentService := services.NewDailyActivityEnrolmentService(testutil.TestDB)

	user1 := models.User{FirstName: "EnrolUser1", LoginEmail: "enrol1@user.com", Active: true}
	testutil.TestDB.Create(&user1)
	user2 := models.User{FirstName: "EnrolUser2", LoginEmail: "enrol2@user.com", Active: true}
	testutil.TestDB.Create(&user2)
	activity1 := models.DailyActivity{Name: "Activity 1", Description: "Desc 1"}
	testutil.TestDB.Create(&activity1)
	activity2 := models.DailyActivity{Name: "Activity 2", Description: "Desc 2"}
	testutil.TestDB.Create(&activity2)

	enrolDTO := ports.EnrolmentInput{
		UserID:          user1.ID,
		DailyActivityID: activity1.ID,
	}
	enrolment, err := enrolmentService.EnrolUser(context.Background(), enrolDTO)
	assert.NoError(t, err)
	assert.NotNil(t, enrolment)

	_, err = enrolmentService.EnrolUser(context.Background(), enrolDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrAlreadyEnrolled, err)

	testutil.TestDB.Create(&models.DailyActivityEnrolment{UserID: user2.ID, DailyActivityID: activity1.ID})
	testutil.TestDB.Create(&models.DailyActivityEnrolment{UserID: user1.ID, DailyActivityID: activity2.ID})

	activity1Enrolments, err := enrolmentService.GetEnrolmentsForActivity(context.Background(), activity1.ID)
	assert.NoError(t, err)
	assert.Len(t, activity1Enrolments, 2)
	assert.Equal(t, "EnrolUser1", activity1Enrolments[0].User.FirstName)

	user1Enrolments, err := enrolmentService.GetEnrolmentsForUser(context.Background(), user1.ID)
	assert.NoError(t, err)
	assert.Len(t, user1Enrolments, 2)
	assert.Equal(t, "Activity 1", user1Enrolments[0].DailyActivity.Name)

	err = enrolmentService.WithdrawUser(context.Background(), activity1.ID, user1.ID)
	assert.NoError(t, err)

	activity1Enrolments, _ = enrolmentService.GetEnrolmentsForActivity(context.Background(), activity1.ID)
	assert.Len(t, activity1Enrolments, 1)

	err = enrolmentService.WithdrawUser(context.Background(), activity1.ID, user1.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrEnrolmentNotFound, err)
}
