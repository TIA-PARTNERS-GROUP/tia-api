package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateUserSkillInput struct {
	UserID           uint                        `json:"user_id" validate:"required"`
	SkillID          uint                        `json:"skill_id" validate:"required"`
	ProficiencyLevel models.UserSkillProficiency `json:"proficiency_level" validate:"required,oneof=beginner intermediate advanced expert"`
}

type UpdateUserSkillInput struct {
	ProficiencyLevel *models.UserSkillProficiency `json:"proficiency_level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced expert"`
}

type UserSkillResponse struct {
	SkillID          uint                        `json:"skill_id"`
	UserID           uint                        `json:"user_id"`
	ProficiencyLevel models.UserSkillProficiency `json:"proficiency_level"`
	CreatedAt        time.Time                   `json:"created_at"`
	Skill            SkillResponse               `json:"skill"`
	User             UserResponse                `json:"user"`
}

type UserSkillsResponse struct {
	Skills []UserSkillResponse `json:"skills"`
	Count  int                 `json:"count"`
}

func MapToUserSkillResponse(us *models.UserSkill) UserSkillResponse {
	return UserSkillResponse{
		SkillID:          us.SkillID,
		UserID:           us.UserID,
		ProficiencyLevel: us.ProficiencyLevel,
		CreatedAt:        us.CreatedAt,
		Skill:            MapSkillToResponse(&us.Skill),
		User:             MapUserToResponse(&us.User),
	}
}

func MapToUserSkillsResponse(userSkills []models.UserSkill) UserSkillsResponse {
	skills := make([]UserSkillResponse, len(userSkills))
	for i, userSkill := range userSkills {
		skills[i] = MapToUserSkillResponse(&userSkill)
	}

	return UserSkillsResponse{
		Skills: skills,
		Count:  len(skills),
	}
}
