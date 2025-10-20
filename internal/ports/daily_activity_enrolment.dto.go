package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type EnrolmentInput struct {
	UserID          uint `json:"user_id" validate:"required"`
	DailyActivityID uint `json:"daily_activity_id" validate:"required"`
}

type UserEnrolmentResponse struct {
	UserID        uint                  `json:"user_id"`
	DailyActivity DailyActivityResponse `json:"daily_activity"`
}

type ActivityEnrolmentResponse struct {
	DailyActivityID uint         `json:"daily_activity_id"`
	User            UserResponse `json:"user"`
}

func MapToUserEnrolmentResponse(enrolment *models.DailyActivityEnrolment) UserEnrolmentResponse {
	return UserEnrolmentResponse{
		UserID:        enrolment.UserID,
		DailyActivity: MapDailyActivityToResponse(&enrolment.DailyActivity),
	}
}

func MapToActivityEnrolmentResponse(enrolment *models.DailyActivityEnrolment) ActivityEnrolmentResponse {
	return ActivityEnrolmentResponse{
		DailyActivityID: enrolment.DailyActivityID,
		User:            MapUserToResponse(&enrolment.User),
	}
}
