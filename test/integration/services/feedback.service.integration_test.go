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
func TestFeedbackService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	feedbackService := services.NewFeedbackService(testutil.TestDB)
	createDTO := ports.CreateFeedbackInput{
		Name:    "John Doe",
		Email:   "john@example.com",
		Content: "This is a test feedback message.",
	}
	createdFeedback, err := feedbackService.CreateFeedback(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdFeedback)
	assert.NotZero(t, createdFeedback.ID)
	assert.Equal(t, "John Doe", createdFeedback.Name)
	fetchedFeedback, err := feedbackService.GetFeedbackByID(context.Background(), createdFeedback.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedFeedback)
	assert.Equal(t, "john@example.com", fetchedFeedback.Email)
}
func TestFeedbackService_Integration_GetAllFeedback(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	feedbackService := services.NewFeedbackService(testutil.TestDB)
	testutil.TestDB.Create(&models.Feedback{Name: "Feedback 1", Email: "1@test.com", Content: "..."})
	testutil.TestDB.Create(&models.Feedback{Name: "Feedback 2", Email: "2@test.com", Content: "..."})
	feedbacks, err := feedbackService.GetAllFeedback(context.Background())
	assert.NoError(t, err)
	assert.Len(t, feedbacks, 2)
}
func TestFeedbackService_Integration_DeleteFeedback(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	feedbackService := services.NewFeedbackService(testutil.TestDB)
	feedback := models.Feedback{Name: "ToDelete", Email: "delete@me.com", Content: "..."}
	testutil.TestDB.Create(&feedback)
	t.Run("Success - Delete Existing", func(t *testing.T) {
		err := feedbackService.DeleteFeedback(context.Background(), feedback.ID)
		assert.NoError(t, err)
	})
	t.Run("Failure - Delete Non-Existent", func(t *testing.T) {
		err := feedbackService.DeleteFeedback(context.Background(), 9999)
		assert.Error(t, err)
		assert.Equal(t, ports.ErrFeedbackNotFound, err)
	})
}
