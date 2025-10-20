package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateSubscriptionInput struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	ValidDays   *int    `json:"valid_days" validate:"omitempty,gt=0"`
	ValidMonths *int    `json:"valid_months" validate:"omitempty,gt=0"`
}

type UserSubscribeInput struct {
	UserID         uint `json:"user_id" validate:"required"`
	SubscriptionID uint `json:"subscription_id" validate:"required"`
}

type SubscriptionResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ValidDays   *int    `json:"valid_days,omitempty"`
	ValidMonths *int    `json:"valid_months,omitempty"`
}

func MapSubscriptionToResponse(sub *models.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		Name:        sub.Name,
		Price:       sub.Price,
		ValidDays:   sub.ValidDays,
		ValidMonths: sub.ValidMonths,
	}
}
