package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateDailyActivityInput struct {
	Name        string `json:"name" validate:"required,min=2,max=60"`
	Description string `json:"description" validate:"required"`
}

type UpdateDailyActivityInput struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=60"`
	Description *string `json:"description"`
}

type EnrolInActivityInput struct {
	UserID          uint `json:"user_id" validate:"required"`
	DailyActivityID uint `json:"daily_activity_id" validate:"required"`
}

type DailyActivityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func MapDailyActivityToResponse(da *models.DailyActivity) DailyActivityResponse {
	return DailyActivityResponse{
		ID:          da.ID,
		Name:        da.Name,
		Description: da.Description,
	}
}
