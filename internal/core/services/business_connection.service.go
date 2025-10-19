package services

import (
	"context"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	"gorm.io/gorm"
)

type BusinessConnectionService struct {
	db *gorm.DB
}

func NewBusinessConnectionService(db *gorm.DB) *BusinessConnectionService {
	return &BusinessConnectionService{db: db}
}

func (s *BusinessConnectionService) CreateBusinessConnection(ctx context.Context, data ports.CreateBusinessConnectionInput) (*models.BusinessConnection, error) {
	var initiatingBusiness models.Business
	if err := s.db.WithContext(ctx).First(&initiatingBusiness, data.InitiatingBusinessID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessNotFound
		}
		return nil, ports.ErrDatabase
	}

	var receivingBusiness models.Business
	if err := s.db.WithContext(ctx).First(&receivingBusiness, data.ReceivingBusinessID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessNotFound
		}
		return nil, ports.ErrDatabase
	}

	var user models.User
	if err := s.db.WithContext(ctx).First(&user, data.InitiatedByUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserNotFound
		}
		return nil, ports.ErrDatabase
	}

	if data.InitiatingBusinessID == data.ReceivingBusinessID {
		return nil, ports.ErrCannotConnectToSelf
	}

	var existingConnection models.BusinessConnection
	err := s.db.WithContext(ctx).
		Where("initiating_business_id = ? AND receiving_business_id = ? AND connection_type = ?",
			data.InitiatingBusinessID, data.ReceivingBusinessID, data.ConnectionType).
		First(&existingConnection).Error

	if err == nil {
		return nil, ports.ErrBusinessConnectionAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, ports.ErrDatabase
	}

	businessConnection := models.BusinessConnection{
		InitiatingBusinessID: data.InitiatingBusinessID,
		ReceivingBusinessID:  data.ReceivingBusinessID,
		ConnectionType:       data.ConnectionType,
		Status:               models.ConnectionStatusPending,
		InitiatedByUserID:    data.InitiatedByUserID,
		Notes:                data.Notes,
	}

	if err := s.db.WithContext(ctx).Create(&businessConnection).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		First(&businessConnection, businessConnection.ID).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &businessConnection, nil
}

func (s *BusinessConnectionService) GetBusinessConnection(ctx context.Context, id uint) (*models.BusinessConnection, error) {
	var businessConnection models.BusinessConnection
	err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		First(&businessConnection, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessConnectionNotFound
		}
		return nil, ports.ErrDatabase
	}

	return &businessConnection, nil
}

func (s *BusinessConnectionService) GetBusinessConnections(ctx context.Context, businessID uint, connectionType *models.BusinessConnectionType, status *models.BusinessConnectionStatus) ([]models.BusinessConnection, error) {
	var businessConnections []models.BusinessConnection
	query := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		Where("initiating_business_id = ? OR receiving_business_id = ?", businessID, businessID)

	if connectionType != nil {
		query = query.Where("connection_type = ?", *connectionType)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at desc").Find(&businessConnections).Error
	if err != nil {
		return nil, ports.ErrDatabase
	}

	return businessConnections, nil
}

func (s *BusinessConnectionService) UpdateBusinessConnection(ctx context.Context, id uint, data ports.UpdateBusinessConnectionInput) (*models.BusinessConnection, error) {
	var businessConnection models.BusinessConnection
	err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&businessConnection).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrBusinessConnectionNotFound
		}
		return nil, ports.ErrDatabase
	}

	updates := make(map[string]interface{})
	if data.ConnectionType != nil {
		updates["connection_type"] = *data.ConnectionType
	}
	if data.Status != nil {
		updates["status"] = *data.Status
	}
	if data.Notes != nil {
		updates["notes"] = *data.Notes
	}

	if len(updates) == 0 {
		return nil, ports.ErrNoUpdateData
	}

	if err := s.db.WithContext(ctx).
		Model(&businessConnection).
		Updates(updates).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		First(&businessConnection, id).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &businessConnection, nil
}

func (s *BusinessConnectionService) DeleteBusinessConnection(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).
		Delete(&models.BusinessConnection{}, id)

	if result.Error != nil {
		return ports.ErrDatabase
	}

	if result.RowsAffected == 0 {
		return ports.ErrBusinessConnectionNotFound
	}

	return nil
}

func (s *BusinessConnectionService) AcceptBusinessConnection(ctx context.Context, id uint) (*models.BusinessConnection, error) {
	var businessConnection models.BusinessConnection
	err := s.db.WithContext(ctx).
		Where("id = ? AND status = ?", id, models.ConnectionStatusPending).
		First(&businessConnection).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrConnectionNotPending
		}
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Model(&businessConnection).
		Update("status", models.ConnectionStatusActive).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		First(&businessConnection, id).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &businessConnection, nil
}

func (s *BusinessConnectionService) RejectBusinessConnection(ctx context.Context, id uint) (*models.BusinessConnection, error) {
	var businessConnection models.BusinessConnection
	err := s.db.WithContext(ctx).
		Where("id = ? AND status = ?", id, models.ConnectionStatusPending).
		First(&businessConnection).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrConnectionNotPending
		}
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Model(&businessConnection).
		Update("status", models.ConnectionStatusRejected).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	if err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		First(&businessConnection, id).Error; err != nil {
		return nil, ports.ErrDatabase
	}

	return &businessConnection, nil
}

func (s *BusinessConnectionService) GetPendingConnections(ctx context.Context, businessID uint) ([]models.BusinessConnection, error) {
	var businessConnections []models.BusinessConnection
	err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		Where("receiving_business_id = ? AND status = ?", businessID, models.ConnectionStatusPending).
		Order("created_at desc").
		Find(&businessConnections).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}

	return businessConnections, nil
}

func (s *BusinessConnectionService) GetActiveConnections(ctx context.Context, businessID uint) ([]models.BusinessConnection, error) {
	var businessConnections []models.BusinessConnection
	err := s.db.WithContext(ctx).
		Preload("InitiatingBusiness").
		Preload("ReceivingBusiness").
		Preload("InitiatedByUser").
		Where("(initiating_business_id = ? OR receiving_business_id = ?) AND status = ?",
			businessID, businessID, models.ConnectionStatusActive).
		Order("created_at desc").
		Find(&businessConnections).Error

	if err != nil {
		return nil, ports.ErrDatabase
	}

	return businessConnections, nil
}
