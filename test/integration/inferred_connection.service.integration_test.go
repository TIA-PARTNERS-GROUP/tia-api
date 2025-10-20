package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/stretchr/testify/assert"
)

func TestInferredConnection_Integration_CreateAndGet(t *testing.T) {
	cleanupTestDB(t)
	icService := services.NewInferredConnectionService(testDB)

	user := models.User{FirstName: "SourceUser", LoginEmail: "source@user.com", Active: true}
	testDB.Create(&user)
	skill := models.Skill{Name: "TargetSkill", Category: "ML", Active: true}
	testDB.Model(&models.Skill{}).Create(map[string]interface{}{
		"Name": "TargetSkill", "Category": "ML", "Active": true,
	})

	createDTO := ports.CreateInferredConnectionInput{
		SourceEntityType: "user",
		SourceEntityID:   user.ID,
		TargetEntityType: "skill",
		TargetEntityID:   skill.ID,
		ConnectionType:   "Recommended_Skill",
		ConfidenceScore:  0.95,
		ModelVersion:     "v1.2.0",
	}
	createdIC, err := icService.CreateInferredConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdIC)
	assert.NotZero(t, createdIC.ID)
	assert.Equal(t, 0.95, createdIC.ConfidenceScore)

	connections, err := icService.GetConnectionsForSource(context.Background(), "user", user.ID)
	assert.NoError(t, err)
	assert.Len(t, connections, 1)
	assert.Equal(t, "Recommended_Skill", connections[0].ConnectionType)
	assert.Equal(t, skill.ID, connections[0].TargetEntityID)
}
