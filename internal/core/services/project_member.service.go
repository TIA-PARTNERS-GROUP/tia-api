package services
import (
	"context"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type ProjectMemberService struct {
	db *gorm.DB
}
func NewProjectMemberService(db *gorm.DB) *ProjectMemberService {
	return &ProjectMemberService{db: db}
}
func (s *ProjectMemberService) AddProjectMember(ctx context.Context, data ports.AddProjectMemberInput) (*models.ProjectMember, error) {
	var project models.Project
	if err := s.db.WithContext(ctx).First(&project, data.ProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectNotFound
		}
		return nil, ports.ErrDatabase
	}
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}
	var existingProjectMember models.ProjectMember
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", data.ProjectID, data.UserID).
		First(&existingProjectMember).Error
	if err == nil {
		return nil, ports.ErrProjectMemberAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, ports.ErrDatabase
	}
	projectMember := models.ProjectMember{
		ProjectID: data.ProjectID,
		UserID:    data.UserID,
		Role:      data.Role,
	}
	if err := s.db.WithContext(ctx).Create(&projectMember).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("User").
		First(&projectMember, "project_id = ? AND user_id = ?", data.ProjectID, data.UserID).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &projectMember, nil
}
func (s *ProjectMemberService) GetProjectMember(ctx context.Context, projectID, userID uint) (*models.ProjectMember, error) {
	var projectMember models.ProjectMember
	err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("User").
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&projectMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectMemberNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &projectMember, nil
}
func (s *ProjectMemberService) GetProjectMembers(ctx context.Context, projectID uint) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("User").
		Where("project_id = ?", projectID).
		Order("role desc, joined_at asc").
		Find(&projectMembers).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return projectMembers, nil
}
func (s *ProjectMemberService) GetProjectsByUser(ctx context.Context, userID uint, role *models.ProjectMemberRole) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	query := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Where("user_id = ?", userID)
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	err := query.Order("joined_at desc").Find(&projectMembers).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return projectMembers, nil
}
func (s *ProjectMemberService) UpdateProjectMemberRole(ctx context.Context, projectID, userID uint, data ports.UpdateProjectMemberRoleInput) (*models.ProjectMember, error) {
	var projectMember models.ProjectMember
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&projectMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectMemberNotFound
		}
		return nil, ports.ErrDatabase
	}
	if projectMember.Role == models.ProjectMemberRoleManager {
	}
	if err := s.db.WithContext(ctx).
		Model(&projectMember).
		Update("role", data.Role).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("User").
		First(&projectMember, "project_id = ? AND user_id = ?", projectID, userID).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &projectMember, nil
}
func (s *ProjectMemberService) RemoveProjectMember(ctx context.Context, projectID, userID uint) error {
	var projectMember models.ProjectMember
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&projectMember).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ports.ErrProjectMemberNotFound
		}
		return ports.ErrDatabase
	}
	if projectMember.Role == models.ProjectMemberRoleManager {
		return ports.ErrCannotRemoveManager
	}
	result := s.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{})
	if result.Error != nil {
		return ports.ErrDatabase
	}
	return nil
}
func (s *ProjectMemberService) GetMembersByRole(ctx context.Context, projectID uint, role models.ProjectMemberRole) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	err := s.db.WithContext(ctx).
		Preload("User").
		Where("project_id = ? AND role = ?", projectID, role).
		Order("joined_at asc").
		Find(&projectMembers).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return projectMembers, nil
}
func (s *ProjectMemberService) IsUserProjectMember(ctx context.Context, projectID, userID uint) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	if err != nil {
		return false, ports.ErrDatabase
	}
	return count > 0, nil
}
