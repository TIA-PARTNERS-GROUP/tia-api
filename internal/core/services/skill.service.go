package services

import (
	"context"
	"errors"
	"strings"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type SkillService struct {
	db *gorm.DB
}

func NewSkillService(db *gorm.DB) *SkillService {
	return &SkillService{db: db}
}
func (s *SkillService) GetSkills(ctx context.Context, filters ports.SkillsFilter) ([]models.Skill, error) {
	var skills []models.Skill
	query := s.db.WithContext(ctx).Order("name asc")
	if filters.Category != nil {
		query = query.Where("category LIKE ?", "%"+*filters.Category+"%")
	}
	if filters.Active != nil {
		query = query.Where("active = ?", *filters.Active)
	}
	if filters.Search != nil {
		searchQuery := "%" + *filters.Search + "%"
		query = query.Where("name LIKE ? OR category LIKE ? OR description LIKE ?", searchQuery, searchQuery, searchQuery)
	}
	if err := query.Find(&skills).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return skills, nil
}
func (s *SkillService) GetSkillByID(ctx context.Context, id uint) (*models.Skill, error) {
	var skill models.Skill
	if err := s.db.WithContext(ctx).First(&skill, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrSkillNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &skill, nil
}
func (s *SkillService) CreateSkill(ctx context.Context, data ports.CreateSkillInput) (*models.Skill, error) {
	active := true
	if data.Active != nil {
		active = *data.Active
	}
	skill := models.Skill{
		Category:    data.Category,
		Name:        data.Name,
		Description: data.Description,
		Active:      active,
	}
	if err := s.db.WithContext(ctx).Create(&skill).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrSkillNameExists
		}
		return nil, ports.ErrDatabase
	}
	return &skill, nil
}
func (s *SkillService) UpdateSkill(ctx context.Context, id uint, data ports.UpdateSkillInput) (*models.Skill, error) {
	skill, err := s.GetSkillByID(ctx, id)
	if err != nil {
		return nil, err
	}
	updateData := make(map[string]interface{})
	if data.Category != nil {
		updateData["category"] = *data.Category
	}
	if data.Name != nil {
		updateData["name"] = *data.Name
	}
	if data.Description != nil {
		updateData["description"] = *data.Description
	}
	if data.Active != nil {
		updateData["active"] = *data.Active
	}
	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}
	if err := s.db.WithContext(ctx).Model(skill).Updates(updateData).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrSkillNameExists
		}
		return nil, ports.ErrDatabase
	}
	return skill, nil
}
func (s *SkillService) DeleteSkill(ctx context.Context, id uint) error {
	var userSkillCount int64
	s.db.Model(&models.UserSkill{}).Where("skill_id = ?", id).Count(&userSkillCount)
	var projectSkillCount int64
	s.db.Model(&models.ProjectSkill{}).Where("skill_id = ?", id).Count(&projectSkillCount)
	if userSkillCount > 0 || projectSkillCount > 0 {
		return ports.ErrSkillInUse
	}
	result := s.db.WithContext(ctx).Delete(&models.Skill{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrSkillNotFound
	}
	return nil
}
func (s *SkillService) ToggleSkillStatus(ctx context.Context, id uint) (*models.Skill, error) {
	skill, err := s.GetSkillByID(ctx, id)
	if err != nil {
		return nil, err
	}
	newStatus := !skill.Active
	if err := s.db.WithContext(ctx).Model(skill).Update("active", newStatus).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	skill.Active = newStatus
	return skill, nil
}
