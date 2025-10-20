package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"gorm.io/datatypes"
)

type SetUserConfigInput struct {
	ConfigType string         `json:"config_type" validate:"required"`
	Config     datatypes.JSON `json:"config" validate:"required"`
}

type UserConfigResponse struct {
	UserID     uint           `json:"user_id"`
	ConfigType string         `json:"config_type"`
	Config     datatypes.JSON `json:"config"`
}

func MapUserConfigToResponse(uc *models.UserConfig) UserConfigResponse {
	return UserConfigResponse{
		UserID:     uc.UserID,
		ConfigType: uc.ConfigType,
		Config:     uc.Config,
	}
}
