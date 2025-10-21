package services
import (
	"context"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type UserSkillService struct {
	db *gorm.DB
}
func NewUserSkillService(db *gorm.DB) *UserSkillService {
	return &UserSkillService{db: db}
}
func (s *UserSkillService) AddUserSkill(ctx context.Context, data ports.CreateUserSkillInput) (*models.UserSkill, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserNotFound
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
	var existingUserSkill models.UserSkill
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND skill_id = ?", data.UserID, data.SkillID).
		First(&existingUserSkill).Error
	if err == nil {
		return nil, ports.ErrUserSkillAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, ports.ErrDatabase
	}
	userSkill := models.UserSkill{
		UserID:           data.UserID,
		SkillID:          data.SkillID,
		ProficiencyLevel: data.ProficiencyLevel,
	}
	if err := s.db.WithContext(ctx).Create(&userSkill).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Skill").
		Preload("User").
		First(&userSkill, "user_id = ? AND skill_id = ?", data.UserID, data.SkillID).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &userSkill, nil
}
func (s *UserSkillService) GetUserSkill(ctx context.Context, userID, skillID uint) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	err := s.db.WithContext(ctx).
		Preload("Skill").
		Preload("User").
		Where("user_id = ? AND skill_id = ?", userID, skillID).
		First(&userSkill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserSkillNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &userSkill, nil
}
func (s *UserSkillService) GetUserSkills(ctx context.Context, userID uint) ([]models.UserSkill, error) {
	var userSkills []models.UserSkill
	err := s.db.WithContext(ctx).
		Preload("Skill").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&userSkills).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return userSkills, nil
}
func (s *UserSkillService) UpdateUserSkill(ctx context.Context, userID, skillID uint, data ports.UpdateUserSkillInput) (*models.UserSkill, error) {
	var userSkill models.UserSkill
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND skill_id = ?", userID, skillID).
		First(&userSkill).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserSkillNotFound
		}
		return nil, ports.ErrDatabase
	}
	updates := make(map[string]interface{})
	if data.ProficiencyLevel != nil {
		updates["proficiency_level"] = *data.ProficiencyLevel
	}
	if len(updates) == 0 {
		return nil, ports.ErrNoUpdateData
	}
	if err := s.db.WithContext(ctx).
		Model(&userSkill).
		Updates(updates).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Skill").
		Preload("User").
		First(&userSkill, "user_id = ? AND skill_id = ?", userID, skillID).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &userSkill, nil
}
func (s *UserSkillService) RemoveUserSkill(ctx context.Context, userID, skillID uint) error {
	result := s.db.WithContext(ctx).
		Where("user_id = ? AND skill_id = ?", userID, skillID).
		Delete(&models.UserSkill{})
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrUserSkillNotFound
	}
	return nil
}
func (s *UserSkillService) GetUsersBySkill(ctx context.Context, skillID uint, proficiency *models.UserSkillProficiency) ([]models.UserSkill, error) {
	var userSkills []models.UserSkill
	query := s.db.WithContext(ctx).
		Preload("User").
		Preload("Skill").
		Where("skill_id = ?", skillID)
	if proficiency != nil {
		query = query.Where("proficiency_level = ?", *proficiency)
	}
	err := query.Order("created_at desc").Find(&userSkills).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return userSkills, nil
}
