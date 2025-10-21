package services
import (
	"context"
	"strings"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)
type ProjectRegionService struct {
	db *gorm.DB
}
func NewProjectRegionService(db *gorm.DB) *ProjectRegionService {
	return &ProjectRegionService{db: db}
}
func (s *ProjectRegionService) AddRegionToProject(ctx context.Context, data ports.AddProjectRegionInput) (*models.ProjectRegion, error) {
	association := models.ProjectRegion{
		ProjectID: data.ProjectID,
		RegionID:  data.RegionID,
	}
	if err := s.db.WithContext(ctx).Create(&association).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ports.ErrRegionAlreadyAdded
		}
		if strings.Contains(err.Error(), "FOREIGN KEY") {
			return nil, ports.ErrProjectOrRegionNotFound
		}
		return nil, ports.ErrDatabase
	}
	return &association, nil
}
func (s *ProjectRegionService) RemoveRegionFromProject(ctx context.Context, projectID uint, regionID string) error {
	result := s.db.WithContext(ctx).Delete(&models.ProjectRegion{}, "project_id = ? AND region_id = ?", projectID, regionID)
	if result.Error != nil {
		return ports.ErrDatabase
	}
	if result.RowsAffected == 0 {
		return ports.ErrProjectRegionNotFound
	}
	return nil
}
func (s *ProjectRegionService) GetRegionsForProject(ctx context.Context, projectID uint) ([]models.ProjectRegion, error) {
	var projectRegions []models.ProjectRegion
	err := s.db.WithContext(ctx).
		Preload("Region").
		Where("project_id = ?", projectID).
		Order("region_id asc").
		Find(&projectRegions).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return projectRegions, nil
}
