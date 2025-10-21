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
func (s *BusinessService) UpdateBusiness(ctx context.Context, id uint, authUserID uint, input ports.UpdateBusinessInput) (*models.Business, error) {
	
	var business models.Business
	if err := s.db.WithContext(ctx).First(&business, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrBusinessNotFound
		}
		return nil, ports.ErrDatabase
	}
	
	if business.OperatorUserID != authUserID {
		
		
		return nil, ports.ErrForbidden
	}
	
	
	updated := false
	if input.Name != nil {
		business.Name = *input.Name
		updated = true
	}
	if input.Tagline != nil {
		business.Tagline = input.Tagline 
		updated = true
	}
	if input.Website != nil {
		business.Website = input.Website
		updated = true
	}
	if input.ContactName != nil {
		business.ContactName = input.ContactName
		updated = true
	}
	if input.ContactPhoneNo != nil {
		business.ContactPhoneNo = input.ContactPhoneNo
		updated = true
	}
	if input.ContactEmail != nil {
		business.ContactEmail = input.ContactEmail
		updated = true
	}
	if input.Description != nil {
		business.Description = input.Description
		updated = true
	}
	if input.Address != nil {
		business.Address = input.Address
		updated = true
	}
	if input.City != nil {
		business.City = input.City
		updated = true
	}
	if input.State != nil {
		business.State = input.State
		updated = true
	}
	if input.Country != nil {
		business.Country = input.Country
		updated = true
	}
	if input.PostalCode != nil {
		business.PostalCode = input.PostalCode
		updated = true
	}
	if input.Value != nil {
		business.Value = input.Value
		updated = true
	}
	if input.BusinessType != nil {
		business.BusinessType = *input.BusinessType
		updated = true
	}
	if input.BusinessCategory != nil {
		business.BusinessCategory = *input.BusinessCategory
		updated = true
	}
	if input.BusinessPhase != nil {
		business.BusinessPhase = *input.BusinessPhase
		updated = true
	}
	if input.Active != nil {
		business.Active = *input.Active
		updated = true
	}
	
	if !updated {
		return nil, ports.ErrNoUpdateData
	}
	
	if err := s.db.WithContext(ctx).Save(&business).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	
	return s.GetBusinessByID(ctx, id)
}
func (s *BusinessService) DeleteBusiness(ctx context.Context, id uint, authUserID uint) error { 
	
	var business models.Business
	if err := s.db.WithContext(ctx).Select("operator_user_id").First(&business, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ports.ErrBusinessNotFound 
		}
		return ports.ErrDatabase 
	}
	if business.OperatorUserID != authUserID {
		
		return ports.ErrForbidden 
	}
	
	
	var projectCount int64
	if err := s.db.Model(&models.Project{}).Where("business_id = ?", id).Count(&projectCount).Error; err != nil {
		return ports.ErrDatabase 
	}
	var pubCount int64
	if err := s.db.Model(&models.Publication{}).Where("business_id = ?", id).Count(&pubCount).Error; err != nil {
		return ports.ErrDatabase 
	}
	if projectCount > 0 || pubCount > 0 {
		return ports.ErrBusinessInUse
	}
	
	result := s.db.WithContext(ctx).Delete(&models.Business{}, id)
	if result.Error != nil {
		return ports.ErrDatabase
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
