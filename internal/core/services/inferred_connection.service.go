package services

import (
	"context"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type InferredConnectionService struct {
	db *gorm.DB
}

func NewInferredConnectionService(db *gorm.DB) *InferredConnectionService {
	return &InferredConnectionService{db: db}
}

func (s *InferredConnectionService) CreateInferredConnection(ctx context.Context, data ports.CreateInferredConnectionInput) (*models.InferredConnection, error) {
	connection := models.InferredConnection{
		SourceEntityType: data.SourceEntityType,
		SourceEntityID:   data.SourceEntityID,
		TargetEntityType: data.TargetEntityType,
		TargetEntityID:   data.TargetEntityID,
		ConnectionType:   data.ConnectionType,
		ConfidenceScore:  data.ConfidenceScore,
		ModelVersion:     data.ModelVersion,
	}
	if err := s.db.WithContext(ctx).Create(&connection).Error; err != nil {
		return nil, ports.ErrDatabase
	}
	return &connection, nil
}

func (s *InferredConnectionService) GetConnectionsForSource(ctx context.Context, entityType string, entityID uint) ([]models.InferredConnection, error) {
	var connections []models.InferredConnection
	err := s.db.WithContext(ctx).
		Where("source_entity_type = ? AND source_entity_id = ?", entityType, entityID).
		Order("confidence_score desc").
		Find(&connections).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}
	return connections, nil
}
