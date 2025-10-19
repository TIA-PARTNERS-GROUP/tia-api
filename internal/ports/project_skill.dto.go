package ports

import (
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateProjectSkillInput struct {
	ProjectID  uint                          `json:"project_id" validate:"required"`
	SkillID    uint                          `json:"skill_id" validate:"required"`
	Importance models.ProjectSkillImportance `json:"importance" validate:"required,oneof=required preferred optional"`
}

type UpdateProjectSkillInput struct {
	Importance *models.ProjectSkillImportance `json:"importance,omitempty" validate:"omitempty,oneof=required preferred optional"`
}

type ProjectSkillResponse struct {
	ProjectID  uint                          `json:"project_id"`
	SkillID    uint                          `json:"skill_id"`
	Importance models.ProjectSkillImportance `json:"importance"`
	Project    ProjectResponse               `json:"project"`
	Skill      SkillResponse                 `json:"skill"`
}

type ProjectSkillsResponse struct {
	Skills []ProjectSkillResponse `json:"skills"`
	Count  int                    `json:"count"`
}

func MapToProjectSkillResponse(ps *models.ProjectSkill) ProjectSkillResponse {
	return ProjectSkillResponse{
		ProjectID:  ps.ProjectID,
		SkillID:    ps.SkillID,
		Importance: ps.Importance,
		Project:    MapToProjectResponse(&ps.Project),
		Skill:      MapSkillToResponse(&ps.Skill),
	}
}

func MapToProjectSkillsResponse(projectSkills []models.ProjectSkill) ProjectSkillsResponse {
	skills := make([]ProjectSkillResponse, len(projectSkills))
	for i, projectSkill := range projectSkills {
		skills[i] = MapToProjectSkillResponse(&projectSkill)
	}

	return ProjectSkillsResponse{
		Skills: skills,
		Count:  len(skills),
	}
}
