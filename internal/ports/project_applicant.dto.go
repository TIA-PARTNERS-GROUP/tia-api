package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type ApplyToProjectInput struct {
	ProjectID uint `json:"project_id" validate:"required"`
	UserID    uint `json:"user_id" validate:"required"`
}

type ProjectApplicantResponse struct {
	ProjectID uint         `json:"project_id"`
	UserID    uint         `json:"user_id"`
	User      UserResponse `json:"user"`
}

func MapProjectApplicantToResponse(pa *models.ProjectApplicant) ProjectApplicantResponse {
	resp := ProjectApplicantResponse{
		ProjectID: pa.ProjectID,
		UserID:    pa.UserID,
	}
	if pa.User.ID != 0 {
		resp.User = MapUserToResponse(&pa.User)
	}
	return resp
}
