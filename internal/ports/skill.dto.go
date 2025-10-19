package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateSkillInput struct {
	Category    string  `json:"category" validate:"required,max=100"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
	Active      *bool   `json:"active"`
}

type UpdateSkillInput struct {
	Category    *string `json:"category" validate:"omitempty,max=100"`
	Name        *string `json:"name" validate:"omitempty,max=100"`
	Description *string `json:"description"`
	Active      *bool   `json:"active"`
}

type SkillsFilter struct {
	Category *string `form:"category"`
	Active   *bool   `form:"active"`
	Search   *string `form:"search"`
}

type SkillResponse struct {
	ID          uint      `json:"id"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

func MapSkillToResponse(skill *models.Skill) SkillResponse {
	return SkillResponse{
		ID:          skill.ID,
		Category:    skill.Category,
		Name:        skill.Name,
		Description: skill.Description,
		Active:      skill.Active,
		CreatedAt:   skill.CreatedAt,
	}
}
