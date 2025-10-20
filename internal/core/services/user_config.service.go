package services

import (
	"context"
	"errors"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserConfigService struct {
	db *gorm.DB
}

func NewUserConfigService(db *gorm.DB) *UserConfigService {
	return &UserConfigService{db: db}
}

func (s *UserConfigService) SetUserConfig(ctx context.Context, userID uint, data ports.SetUserConfigInput) (*models.UserConfig, error) {
	config := models.UserConfig{
		UserID:     userID,
		ConfigType: data.ConfigType,
		Config:     data.Config,
	}

	err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "config_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"config"}),
	}).Create(&config).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}
	return &config, nil
}

func (s *UserConfigService) GetUserConfig(ctx context.Context, userID uint, configType string) (*models.UserConfig, error) {
	var config models.UserConfig
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND config_type = ?", userID, configType).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrUserConfigNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &config, nil
}

func (s *UserConfigService) DeleteUserConfig(ctx context.Context, userID uint, configType string) error {
	result := s.db.WithContext(ctx).
		Delete(&models.UserConfig{}, "user_id = ? AND config_type = ?", userID, configType)

	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrUserConfigNotFound
	}
	return nil
}
