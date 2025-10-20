package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestUserConfigService_Integration_SetAndGet(t *testing.T) {
	cleanupTestDB(t)
	configService := services.NewUserConfigService(testDB)

	user := models.User{FirstName: "Config", LoginEmail: "config@user.com", Active: true}
	testDB.Create(&user)

	dashboardLayout := `{"widgets": ["feed", "stats"], "columns": 2}`
	createDTO := ports.SetUserConfigInput{
		ConfigType: "dashboard_layout",
		Config:     datatypes.JSON(dashboardLayout),
	}
	createdConfig, err := configService.SetUserConfig(context.Background(), user.ID, createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdConfig)
	assert.Equal(t, "dashboard_layout", createdConfig.ConfigType)

	fetchedConfig, err := configService.GetUserConfig(context.Background(), user.ID, "dashboard_layout")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedConfig)

	var configData map[string]interface{}
	err = json.Unmarshal(fetchedConfig.Config, &configData)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), configData["columns"])

	newLayout := `{"widgets": ["feed", "stats", "calendar"], "columns": 3}`
	updateDTO := ports.SetUserConfigInput{
		ConfigType: "dashboard_layout",
		Config:     datatypes.JSON(newLayout),
	}
	updatedConfig, err := configService.SetUserConfig(context.Background(), user.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedConfig)

	fetchedAgain, _ := configService.GetUserConfig(context.Background(), user.ID, "dashboard_layout")
	var newConfigData map[string]interface{}
	json.Unmarshal(fetchedAgain.Config, &newConfigData)
	assert.Equal(t, float64(3), newConfigData["columns"])
}

func TestUserConfigService_Integration_Delete(t *testing.T) {
	cleanupTestDB(t)
	configService := services.NewUserConfigService(testDB)

	user := models.User{FirstName: "ConfigDel", LoginEmail: "configdel@user.com", Active: true}
	testDB.Create(&user)

	config := models.UserConfig{
		UserID:     user.ID,
		ConfigType: "temp_settings",
		Config:     datatypes.JSON(`{"theme": "dark"}`),
	}
	testDB.Create(&config)

	err := configService.DeleteUserConfig(context.Background(), user.ID, "temp_settings")
	assert.NoError(t, err)

	_, err = configService.GetUserConfig(context.Background(), user.ID, "temp_settings")
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserConfigNotFound, err)
}
