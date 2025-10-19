package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateNotificationInput struct {
	SenderUserID      *uint                     `json:"sender_user_id"`
	ReceiverUserID    uint                      `json:"receiver_user_id" validate:"required"`
	NotificationType  models.NotificationType   `json:"notification_type" validate:"required"`
	Title             string                    `json:"title" validate:"required,max=255"`
	Message           string                    `json:"message" validate:"required"`
	RelatedEntityType *models.RelatedEntityType `json:"related_entity_type"`
	RelatedEntityID   *uint                     `json:"related_entity_id"`
	ActionURL         *string                   `json:"action_url" validate:"omitempty,url"`
}

type NotificationResponse struct {
	ID                uint                      `json:"id"`
	NotificationType  models.NotificationType   `json:"notification_type"`
	Title             string                    `json:"title"`
	Message           string                    `json:"message"`
	RelatedEntityType *models.RelatedEntityType `json:"related_entity_type,omitempty"`
	RelatedEntityID   *uint                     `json:"related_entity_id,omitempty"`
	Read              bool                      `json:"read"`
	ActionURL         *string                   `json:"action_url,omitempty"`
	CreatedAt         time.Time                 `json:"created_at"`
	Sender            *UserResponse             `json:"sender,omitempty"`
	Receiver          UserResponse              `json:"receiver"`
}

func MapNotificationToResponse(n *models.Notification) NotificationResponse {
	resp := NotificationResponse{
		ID:                n.ID,
		NotificationType:  n.NotificationType,
		Title:             n.Title,
		Message:           n.Message,
		RelatedEntityType: n.RelatedEntityType,
		RelatedEntityID:   n.RelatedEntityID,
		Read:              n.Read,
		ActionURL:         n.ActionURL,
		CreatedAt:         n.CreatedAt,
	}

	if n.ReceiverUser.ID != 0 {
		resp.Receiver = MapUserToResponse(&n.ReceiverUser)
	}

	if n.SenderUser != nil && n.SenderUser.ID != 0 {
		senderResp := MapUserToResponse(n.SenderUser)
		resp.Sender = &senderResp
	}

	return resp
}
