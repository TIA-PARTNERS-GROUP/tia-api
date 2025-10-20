package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateProjectInput struct {
	ManagedByUserID uint                 `json:"managed_by_user_id" validate:"required"`
	BusinessID      *uint                `json:"business_id"`
	Name            string               `json:"name" validate:"required,min=2,max=100"`
	Description     *string              `json:"description"`
	ProjectStatus   models.ProjectStatus `json:"project_status" validate:"required"`
	StartDate       *time.Time           `json:"start_date"`
	TargetEndDate   *time.Time           `json:"target_end_date"`
	RegionIDs       []string             `json:"region_ids"`
}

type UpdateProjectInput struct {
	ManagedByUserID *uint                 `json:"managed_by_user_id"`
	BusinessID      *uint                 `json:"business_id"`
	Name            *string               `json:"name" validate:"omitempty,min=2,max=100"`
	Description     *string               `json:"description"`
	ProjectStatus   *models.ProjectStatus `json:"project_status"`
	StartDate       *time.Time            `json:"start_date"`
	TargetEndDate   *time.Time            `json:"target_end_date"`
	ActualEndDate   *time.Time            `json:"actual_end_date"`
	RegionIDs       []string              `json:"region_ids"`
}

type AddMemberInput struct {
	UserID uint                     `json:"user_id" validate:"required"`
	Role   models.ProjectMemberRole `json:"role" validate:"required"`
}

type UpdateMemberRoleInput struct {
	Role models.ProjectMemberRole `json:"role" validate:"required"`
}

type ProjectResponse struct {
	ID            uint                    `json:"id"`
	Name          string                  `json:"name"`
	Description   *string                 `json:"description,omitempty"`
	ProjectStatus models.ProjectStatus    `json:"project_status"`
	StartDate     *time.Time              `json:"start_date,omitempty"`
	TargetEndDate *time.Time              `json:"target_end_date,omitempty"`
	ActualEndDate *time.Time              `json:"actual_end_date,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Manager       UserResponse            `json:"manager"`
	Business      *BusinessResponse       `json:"business,omitempty"`
	Members       []ProjectMemberResponse `json:"members"`
	Regions       []RegionResponse        `json:"regions,omitempty"`
}

func MapToProjectResponse(p *models.Project) ProjectResponse {
	resp := ProjectResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		ProjectStatus: p.ProjectStatus,
		StartDate:     p.StartDate,
		TargetEndDate: p.TargetEndDate,
		ActualEndDate: p.ActualEndDate,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.ManagingUser.ID != 0 {
		resp.Manager = MapUserToResponse(&p.ManagingUser)
	}

	if p.Business != nil {
		businessResp := MapBusinessToResponse(p.Business)
		resp.Business = &businessResp
	}

	members := make([]ProjectMemberResponse, len(p.ProjectMembers))
	for i, member := range p.ProjectMembers {
		members[i] = MapToProjectMemberResponse(&member)
	}
	resp.Members = members

	regions := make([]RegionResponse, len(p.ProjectRegions))
	for i, pr := range p.ProjectRegions {
		regions[i] = MapRegionToResponse(&pr.Region)
	}
	resp.Regions = regions

	return resp
}
