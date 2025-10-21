package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"gorm.io/datatypes"
)

type CreateL2EResponseInput struct {
	UserID   uint           `json:"-" validate:"-"`
	Response datatypes.JSON `json:"response" validate:"required"`
}

type L2EResponseResponse struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	Response  datatypes.JSON `json:"response"`
	DateAdded time.Time      `json:"date_added"`
}

func MapL2EResponseToResponse(l2e *models.L2EResponse) L2EResponseResponse {
	return L2EResponseResponse{
		ID:        l2e.ID,
		UserID:    l2e.UserID,
		Response:  l2e.Response,
		DateAdded: l2e.DateAdded,
	}
}
