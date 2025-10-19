package services

import (
	"context"
	"errors"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type BusinessService struct {
	db *gorm.DB
}

func NewBusinessService(db *gorm.DB) *BusinessService {
	return &BusinessService{db: db}
}

func (s *BusinessService) GetBusinesses(ctx context.Context, filters ports.BusinessesFilter) ([]models.Business, error) {
	var businesses []models.Business
	query := s.db.WithContext(ctx).Preload("OperatorUser").Order("name asc")

	if filters.BusinessType != nil {
		query = query.Where("business_type = ?", *filters.BusinessType)
	}
	if filters.BusinessCategory != nil {
		query = query.Where("business_category = ?", *filters.BusinessCategory)
	}
	if filters.BusinessPhase != nil {
		query = query.Where("business_phase = ?", *filters.BusinessPhase)
	}
	if filters.Active != nil {
		query = query.Where("active = ?", *filters.Active)
	}
	if filters.OperatorUserID != nil {
		query = query.Where("operator_user_id = ?", *filters.OperatorUserID)
	}
	if filters.Search != nil {
		searchQuery := "%" + *filters.Search + "%"
		query = query.Where("name LIKE ? OR tagline LIKE ? OR description LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	if err := query.Find(&businesses).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return businesses, nil
}

func (s *BusinessService) GetBusinessByID(ctx context.Context, id uint) (*models.Business, error) {
	var business models.Business
	err := s.db.WithContext(ctx).Preload("OperatorUser").First(&business, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrBusinessNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &business, nil
}

func (s *BusinessService) CreateBusiness(ctx context.Context, data ports.CreateBusinessInput) (*models.Business, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.OperatorUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrOperatorNotFound
		}
		return nil, ports.ErrDatabase
	}

	business := models.Business{
		OperatorUserID:   data.OperatorUserID,
		Name:             data.Name,
		Tagline:          data.Tagline,
		Website:          data.Website,
		ContactName:      data.ContactName,
		ContactPhoneNo:   data.ContactPhoneNo,
		ContactEmail:     data.ContactEmail,
		Description:      data.Description,
		Address:          data.Address,
		City:             data.City,
		State:            data.State,
		Country:          data.Country,
		PostalCode:       data.PostalCode,
		Value:            data.Value,
		BusinessType:     data.BusinessType,
		BusinessCategory: data.BusinessCategory,
		BusinessPhase:    data.BusinessPhase,
		Active:           true,
	}

	if err := s.db.WithContext(ctx).Create(&business).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	business.OperatorUser = user
	return &business, nil
}

func (s *BusinessService) UpdateBusiness(ctx context.Context, id uint, data ports.UpdateBusinessInput) (*models.Business, error) {
	business, err := s.GetBusinessByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})
	if data.Name != nil {
		updateData["name"] = *data.Name
	}
	if data.Tagline != nil {
		updateData["tagline"] = *data.Tagline
	}
	if data.Website != nil {
		updateData["website"] = *data.Website
	}
	if data.ContactName != nil {
		updateData["contact_name"] = *data.ContactName
	}
	if data.ContactPhoneNo != nil {
		updateData["contact_phone_no"] = *data.ContactPhoneNo
	}
	if data.ContactEmail != nil {
		updateData["contact_email"] = *data.ContactEmail
	}
	if data.Description != nil {
		updateData["description"] = *data.Description
	}
	if data.Address != nil {
		updateData["address"] = *data.Address
	}
	if data.City != nil {
		updateData["city"] = *data.City
	}
	if data.State != nil {
		updateData["state"] = *data.State
	}
	if data.Country != nil {
		updateData["country"] = *data.Country
	}
	if data.PostalCode != nil {
		updateData["postal_code"] = *data.PostalCode
	}
	if data.Value != nil {
		updateData["value"] = *data.Value
	}
	if data.BusinessType != nil {
		updateData["business_type"] = *data.BusinessType
	}
	if data.BusinessCategory != nil {
		updateData["business_category"] = *data.BusinessCategory
	}
	if data.BusinessPhase != nil {
		updateData["business_phase"] = *data.BusinessPhase
	}
	if data.Active != nil {
		updateData["active"] = *data.Active
	}

	if len(updateData) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).Model(business).Updates(updateData).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return s.GetBusinessByID(ctx, id)
}

func (s *BusinessService) DeleteBusiness(ctx context.Context, id uint) error {
	var projectCount int64
	s.db.Model(&models.Project{}).Where("business_id = ?", id).Count(&projectCount)

	var pubCount int64
	s.db.Model(&models.Publication{}).Where("business_id = ?", id).Count(&pubCount)

	if projectCount > 0 || pubCount > 0 {
		return ports.ErrBusinessInUse
	}

	result := s.db.WithContext(ctx).Delete(&models.Business{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrBusinessNotFound
	}

	return nil
}

func (s *BusinessService) GetUserBusinesses(ctx context.Context, userID uint) ([]models.Business, error) {
	var businesses []models.Business
	if err := s.db.WithContext(ctx).Where("operator_user_id = ?", userID).Order("name asc").Find(&businesses).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return businesses, nil
}

func (s *BusinessService) ToggleBusinessStatus(ctx context.Context, id uint) (*models.Business, error) {
	business, err := s.GetBusinessByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newStatus := !business.Active
	if err := s.db.WithContext(ctx).Model(business).Update("active", newStatus).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	business.Active = newStatus
	return business, nil
}

func (s *BusinessService) GetBusinessStats(ctx context.Context, id uint) (*ports.BusinessStatsResponse, error) {
	var projectCount, pubCount, tagCount int64

	if err := s.db.Model(&models.Project{}).Where("business_id = ?", id).Count(&projectCount).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.Model(&models.Publication{}).Where("business_id = ?", id).Count(&pubCount).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	if err := s.db.Model(&models.BusinessTag{}).Where("business_id = ?", id).Count(&tagCount).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &ports.BusinessStatsResponse{
		TotalProjects:     projectCount,
		TotalPublications: pubCount,
		TotalTags:         tagCount,
	}, nil
}
