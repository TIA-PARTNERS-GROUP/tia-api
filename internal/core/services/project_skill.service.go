package services

import (
	"context"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type ProjectSkillService struct {
	db *gorm.DB
}

func NewProjectSkillService(db *gorm.DB) *ProjectSkillService {
	return &ProjectSkillService{db: db}
}

func (s *ProjectSkillService) AddProjectSkill(ctx context.Context, data ports.CreateProjectSkillInput) (*models.ProjectSkill, error) {
	var project models.Project
	if err := s.db.WithContext(ctx).First(&project, data.ProjectID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectNotFound
		}
		return nil, ports.ErrDatabase
	}

	var skill models.Skill
	if err := s.db.WithContext(ctx).First(&skill, data.SkillID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrSkillNotFound
		}
		return nil, ports.ErrDatabase
	}

	var existingProjectSkill models.ProjectSkill
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND skill_id = ?", data.ProjectID, data.SkillID).
		First(&existingProjectSkill).Error

	if err == nil {
		return nil, ports.ErrProjectSkillAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, ports.ErrDatabase
	}

	projectSkill := models.ProjectSkill{
		ProjectID:  data.ProjectID,
		SkillID:    data.SkillID,
		Importance: data.Importance,
	}

	if err := s.db.WithContext(ctx).Create(&projectSkill).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("Project.ProjectMembers").
		Preload("Project.ProjectMembers.User").
		Preload("Skill").
		First(&projectSkill, "project_id = ? AND skill_id = ?", data.ProjectID, data.SkillID).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &projectSkill, nil
}

func (s *ProjectSkillService) GetProjectSkill(ctx context.Context, projectID, skillID uint) (*models.ProjectSkill, error) {
	var projectSkill models.ProjectSkill
	err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("Project.ProjectMembers").
		Preload("Project.ProjectMembers.User").
		Preload("Skill").
		Where("project_id = ? AND skill_id = ?", projectID, skillID).
		First(&projectSkill).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectSkillNotFound
		}
		return nil, ports.ErrDatabase
	}

	return &projectSkill, nil
}

func (s *ProjectSkillService) GetProjectSkills(ctx context.Context, projectID uint) ([]models.ProjectSkill, error) {
	var projectSkills []models.ProjectSkill
	err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Skill").
		Where("project_id = ?", projectID).
		Order("importance desc, skill_id asc").
		Find(&projectSkills).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}

	return projectSkills, nil
}

func (s *ProjectSkillService) GetProjectsBySkill(ctx context.Context, skillID uint, importance *models.ProjectSkillImportance) ([]models.ProjectSkill, error) {
	var projectSkills []models.ProjectSkill
	query := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("Skill").
		Where("skill_id = ?", skillID)

	if importance != nil {
		query = query.Where("importance = ?", *importance)
	}

	err := query.Order("importance desc, project_id asc").Find(&projectSkills).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}

	return projectSkills, nil
}

func (s *ProjectSkillService) UpdateProjectSkill(ctx context.Context, projectID, skillID uint, data ports.UpdateProjectSkillInput) (*models.ProjectSkill, error) {
	var projectSkill models.ProjectSkill
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND skill_id = ?", projectID, skillID).
		First(&projectSkill).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrProjectSkillNotFound
		}
		return nil, ports.ErrDatabase
	}

	updates := make(map[string]interface{})
	if data.Importance != nil {
		updates["importance"] = *data.Importance
	}

	if len(updates) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).
		Model(&projectSkill).
		Updates(updates).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("Project").
		Preload("Project.ManagingUser").
		Preload("Project.ProjectMembers").
		Preload("Project.ProjectMembers.User").
		Preload("Skill").
		First(&projectSkill, "project_id = ? AND skill_id = ?", projectID, skillID).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &projectSkill, nil
}

func (s *ProjectSkillService) RemoveProjectSkill(ctx context.Context, projectID, skillID uint) error {
	result := s.db.WithContext(ctx).
		Where("project_id = ? AND skill_id = ?", projectID, skillID).
		Delete(&models.ProjectSkill{})

	if result.Error != nil {
		return ports.ErrDatabase
	}

	if result.RowsAffected == 0 {
		return ports.ErrProjectSkillNotFound
	}

	return nil
}

func (s *ProjectSkillService) GetSkillsByImportance(ctx context.Context, projectID uint, importance models.ProjectSkillImportance) ([]models.ProjectSkill, error) {
	var projectSkills []models.ProjectSkill
	err := s.db.WithContext(ctx).
		Preload("Skill").
		Where("project_id = ? AND importance = ?", projectID, importance).
		Order("skill_id asc").
		Find(&projectSkills).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}

	return projectSkills, nil
}
