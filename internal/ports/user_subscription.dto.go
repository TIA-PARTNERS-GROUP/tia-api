package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateUserSubscriptionInput struct {
	UserID         uint `json:"user_id" validate:"required"`
	SubscriptionID uint `json:"subscription_id" validate:"required"`
	IsTrial        bool `json:"is_trial"`
}

type UserSubscriptionResponse struct {
	ID           uint                 `json:"id"`
	DateFrom     time.Time            `json:"date_from"`
	DateTo       time.Time            `json:"date_to"`
	IsTrial      bool                 `json:"is_trial"`
	User         UserResponse         `json:"user"`
	Subscription SubscriptionResponse `json:"subscription"`
}

func MapUserSubscriptionToResponse(us *models.UserSubscription) UserSubscriptionResponse {
	resp := UserSubscriptionResponse{
		ID:       us.ID,
		DateFrom: us.DateFrom,
		DateTo:   us.DateTo,
		IsTrial:  us.IsTrial,
	}
	if us.User.ID != 0 {
		resp.User = MapUserToResponse(&us.User)
	}
	if us.Subscription.ID != 0 {
		resp.Subscription = MapSubscriptionToResponse(&us.Subscription)
	}
	return resp
}
