package services
import (
	"context"
	"strings"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type BusinessTagService struct {
	db *gorm.DB
}
func NewBusinessTagService(db *gorm.DB) *BusinessTagService {
	return &BusinessTagService{db: db}
}
func (s *BusinessTagService) CreateBusinessTag(ctx context.Context, data ports.CreateBusinessTagInput) (*models.BusinessTag, error) {
	var business models.Business
	if err := s.db.WithContext(ctx).First(&business, data.BusinessID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessNotFound
		}
		return nil, ports.ErrDatabase
	}
	var existingTag models.BusinessTag
	err := s.db.WithContext(ctx).
		Where("business_id = ? AND tag_type = ? AND description = ?",
			data.BusinessID, data.TagType, data.Description).
		First(&existingTag).Error
	if err == nil {
		return nil, ports.ErrBusinessTagAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, ports.ErrDatabase
	}
	businessTag := models.BusinessTag{
		BusinessID:  data.BusinessID,
		TagType:     data.TagType,
		Description: data.Description,
	}
	if err := s.db.WithContext(ctx).Create(&businessTag).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Business").
		First(&businessTag, businessTag.ID).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &businessTag, nil
}
func (s *BusinessTagService) GetBusinessTag(ctx context.Context, id uint) (*models.BusinessTag, error) {
	var businessTag models.BusinessTag
	err := s.db.WithContext(ctx).
		Preload("Business").
		First(&businessTag, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessTagNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &businessTag, nil
}
func (s *BusinessTagService) GetBusinessTags(ctx context.Context, businessID uint, tagType *models.BusinessTagType) ([]models.BusinessTag, error) {
	var businessTags []models.BusinessTag
	query := s.db.WithContext(ctx).
		Preload("Business").
		Where("business_id = ?", businessID)
	if tagType != nil {
		query = query.Where("tag_type = ?", *tagType)
	}
	err := query.Order("created_at desc").Find(&businessTags).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return businessTags, nil
}
func (s *BusinessTagService) UpdateBusinessTag(ctx context.Context, id uint, data ports.UpdateBusinessTagInput) (*models.BusinessTag, error) {
	var businessTag models.BusinessTag
	err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&businessTag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessTagNotFound
		}
		return nil, ports.ErrDatabase
	}
	if data.TagType != nil && data.Description != nil {
		var existingTag models.BusinessTag
		err := s.db.WithContext(ctx).
			Where("business_id = ? AND tag_type = ? AND description = ? AND id != ?",
				businessTag.BusinessID, *data.TagType, *data.Description, id).
			First(&existingTag).Error
		if err == nil {
			return nil, ports.ErrBusinessTagAlreadyExists
		} else if err != gorm.ErrRecordNotFound {
			return nil, ports.ErrDatabase
		}
	}
	updates := make(map[string]interface{})
	if data.TagType != nil {
		updates["tag_type"] = *data.TagType
	}
	if data.Description != nil {
		updates["description"] = *data.Description
	}
	if len(updates) == 0 {
		return nil, ports.ErrNoUpdateData
	}
	if err := s.db.WithContext(ctx).
		Model(&businessTag).
		Updates(updates).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.WithContext(ctx).
		Preload("Business").
		First(&businessTag, id).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &businessTag, nil
}
func (s *BusinessTagService) DeleteBusinessTag(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).
		Delete(&models.BusinessTag{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrBusinessTagNotFound
	}
	return nil
}
func (s *BusinessTagService) GetTagsByType(ctx context.Context, businessID uint, tagType models.BusinessTagType) ([]models.BusinessTag, error) {
	var businessTags []models.BusinessTag
	err := s.db.WithContext(ctx).
		Preload("Business").
		Where("business_id = ? AND tag_type = ?", businessID, tagType).
		Order("created_at desc").
		Find(&businessTags).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return businessTags, nil
}
func (s *BusinessTagService) SearchBusinessTags(ctx context.Context, businessID uint, searchTerm string, tagType *models.BusinessTagType) ([]models.BusinessTag, error) {
	var businessTags []models.BusinessTag
	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	query := s.db.WithContext(ctx).
		Preload("Business").
		Where("business_id = ?", businessID).
		Where("LOWER(description) LIKE ?", searchPattern)
	if tagType != nil {
		query = query.Where("tag_type = ?", *tagType)
	}
	err := query.Order("created_at desc").Find(&businessTags).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return businessTags, nil
}
func (s *BusinessTagService) GetBusinessesByTag(ctx context.Context, tagType models.BusinessTagType, description string) ([]models.BusinessTag, error) {
	var businessTags []models.BusinessTag
	err := s.db.WithContext(ctx).
		Preload("Business").
		Where("tag_type = ? AND description = ?", tagType, description).
		Order("created_at desc").
		Find(&businessTags).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return businessTags, nil
}
