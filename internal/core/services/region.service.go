package services

import (
	"context"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type RegionService struct {
	db *gorm.DB
}

func NewRegionService(db *gorm.DB) *RegionService {
	return &RegionService{db: db}
}

func (s *RegionService) SeedRegions(ctx context.Context) error {
	regions := []models.Region{
		{ID: "QLD", Name: "Queensland"},
		{ID: "NSW", Name: "New South Wales"},
		{ID: "VIC", Name: "Victoria"},
		{ID: "TAS", Name: "Tasmania"},
		{ID: "SA", Name: "South Australia"},
		{ID: "WA", Name: "Western Australia"},
		{ID: "NT", Name: "Northern Territory"},
		{ID: "ACT", Name: "Australian Capital Territory"},
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, region := range regions {
			if err := tx.FirstOrCreate(&region).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *RegionService) GetAllRegions(ctx context.Context) ([]models.Region, error) {
	var regions []models.Region
	if err := s.db.WithContext(ctx).Order("name asc").Find(&regions).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return regions, nil
}
