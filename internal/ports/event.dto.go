package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"gorm.io/datatypes"
)

type CreateEventInput struct {
	EventType string         `json:"event_type" validate:"required,max=100"`
	Payload   datatypes.JSON `json:"payload" validate:"required" swaggertype:"object"`
	UserID    *uint          `json:"user_id"`
}
type EventsFilter struct {
	EventType *string `form:"event_type"`
	UserID    *uint   `form:"user_id"`
}
type EventResponse struct {
	ID        uint           `json:"id"`
	EventType string         `json:"event_type"`
	Payload   datatypes.JSON `json:"payload" swaggertype:"object"`
	Timestamp time.Time      `json:"timestamp"`
	UserID    *uint          `json:"user_id,omitempty"`
	User      *UserResponse  `json:"user,omitempty"`
}

func MapEventToResponse(event *models.Event) EventResponse {
	resp := EventResponse{
		ID:        event.ID,
		EventType: event.EventType,
		Payload:   event.Payload,
		Timestamp: event.Timestamp,
		UserID:    event.UserID,
	}
	if event.User != nil && event.User.ID != 0 {
		userResp := MapUserToResponse(event.User)
		resp.User = &userResp
	}
	return resp
}
