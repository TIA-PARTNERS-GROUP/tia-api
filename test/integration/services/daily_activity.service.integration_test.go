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
func TestDailyActivityService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	activityService := services.NewDailyActivityService(testutil.TestDB)
	createDTO := ports.CreateDailyActivityInput{
		Name:        "Morning Standup",
		Description: "A daily sync meeting for the team.",
	}
	createdActivity, err := activityService.CreateDailyActivity(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdActivity)
	assert.Equal(t, "Morning Standup", createdActivity.Name)
	fetchedActivity, err := activityService.GetDailyActivityByID(context.Background(), createdActivity.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedActivity)
	assert.Equal(t, "Morning Standup", fetchedActivity.Name)
}
func TestDailyActivityService_Integration_Enrolment(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	activityService := services.NewDailyActivityService(testutil.TestDB)
	user := models.User{FirstName: "ActivityUser", LoginEmail: "activity@user.com", Active: true}
	testutil.TestDB.Create(&user)
	activity := models.DailyActivity{Name: "Weekly Review", Description: "Review of the week's progress."}
	testutil.TestDB.Create(&activity)
	enrolDTO := ports.EnrolInActivityInput{
		UserID:          user.ID,
		DailyActivityID: activity.ID,
	}
	enrolment, err := activityService.EnrolUserInActivity(context.Background(), enrolDTO)
	assert.NoError(t, err)
	assert.NotNil(t, enrolment)
	_, err = activityService.EnrolUserInActivity(context.Background(), enrolDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrAlreadyEnrolled, err)
	err = activityService.WithdrawUserFromActivity(context.Background(), activity.ID, user.ID)
	assert.NoError(t, err)
	err = activityService.WithdrawUserFromActivity(context.Background(), activity.ID, user.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrEnrolmentNotFound, err)
}
