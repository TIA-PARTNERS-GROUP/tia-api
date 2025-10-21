package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type AddProjectMemberInput struct {
	ProjectID uint                     `json:"project_id" validate:"required"`
	UserID    uint                     `json:"user_id" validate:"required"`
	Role      models.ProjectMemberRole `json:"role" validate:"required,oneof=manager contributor reviewer"`
}
type UpdateProjectMemberRoleInput struct {
	Role models.ProjectMemberRole `json:"role" validate:"required,oneof=manager contributor reviewer"`
}
type ProjectMemberResponse struct {
	ProjectID uint                     `json:"project_id"`
	UserID    uint                     `json:"user_id"`
	Role      models.ProjectMemberRole `json:"role"`
	JoinedAt  time.Time                `json:"joined_at"`
	Project   ProjectResponse          `json:"project"`
	User      UserResponse             `json:"user"`
}
type ProjectMembersResponse struct {
	Members []ProjectMemberResponse `json:"members"`
	Count   int                     `json:"count"`
}
func MapToProjectMemberResponse(pm *models.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{
		ProjectID: pm.ProjectID,
		UserID:    pm.UserID,
		Role:      pm.Role,
		JoinedAt:  pm.JoinedAt,
		Project:   MapToProjectResponse(&pm.Project),
		User:      MapUserToResponse(&pm.User),
	}
}
func MapToProjectMembersResponse(projectMembers []models.ProjectMember) ProjectMembersResponse {
	members := make([]ProjectMemberResponse, len(projectMembers))
	for i, projectMember := range projectMembers {
		members[i] = MapToProjectMemberResponse(&projectMember)
	}
	return ProjectMembersResponse{
		Members: members,
		Count:   len(members),
	}
}
