package main
import (
	"context"
	"encoding/json"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)
func TestL2EResponseService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	l2eService := services.NewL2EResponseService(testutil.TestDB)
	user := models.User{FirstName: "L2E", LoginEmail: "l2e@user.com", Active: true}
	testutil.TestDB.Create(&user)
	jsonPayload := `{"question_id": 1, "answer": "B"}`
	createDTO := ports.CreateL2EResponseInput{
		UserID:   user.ID,
		Response: datatypes.JSON(jsonPayload),
	}
	createdResponse, err := l2eService.CreateL2EResponse(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdResponse)
	assert.NotZero(t, createdResponse.ID)
	assert.Equal(t, user.ID, createdResponse.UserID)
	responses, err := l2eService.GetL2EResponsesForUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, responses, 1)
	var responseData map[string]interface{}
	err = json.Unmarshal(responses[0].Response, &responseData)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), responseData["question_id"])
	assert.Equal(t, "B", responseData["answer"])
}
