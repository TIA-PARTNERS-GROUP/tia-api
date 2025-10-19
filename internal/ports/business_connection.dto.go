package ports

import (
	"time"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

type CreateBusinessConnectionInput struct {
	InitiatingBusinessID uint                          `json:"initiating_business_id" validate:"required"`
	ReceivingBusinessID  uint                          `json:"receiving_business_id" validate:"required"`
	ConnectionType       models.BusinessConnectionType `json:"connection_type" validate:"required,oneof=Partnership Supplier Client Referral Collaboration"`
	InitiatedByUserID    uint                          `json:"initiated_by_user_id" validate:"required"`
	Notes                *string                       `json:"notes,omitempty"`
}

type UpdateBusinessConnectionInput struct {
	ConnectionType *models.BusinessConnectionType   `json:"connection_type,omitempty" validate:"omitempty,oneof=Partnership Supplier Client Referral Collaboration"`
	Status         *models.BusinessConnectionStatus `json:"status,omitempty" validate:"omitempty,oneof=pending active rejected inactive"`
	Notes          *string                          `json:"notes,omitempty"`
}

type BusinessConnectionResponse struct {
	ID                   uint                            `json:"id"`
	InitiatingBusinessID uint                            `json:"initiating_business_id"`
	ReceivingBusinessID  uint                            `json:"receiving_business_id"`
	ConnectionType       models.BusinessConnectionType   `json:"connection_type"`
	Status               models.BusinessConnectionStatus `json:"status"`
	InitiatedByUserID    uint                            `json:"initiated_by_user_id"`
	Notes                *string                         `json:"notes,omitempty"`
	CreatedAt            time.Time                       `json:"created_at"`
	UpdatedAt            time.Time                       `json:"updated_at"`

	// Relationships
	InitiatingBusiness BusinessResponse `json:"initiating_business"`
	ReceivingBusiness  BusinessResponse `json:"receiving_business"`
	InitiatedByUser    UserResponse     `json:"initiated_by_user"`
}

type BusinessConnectionsResponse struct {
	Connections []BusinessConnectionResponse `json:"connections"`
	Count       int                          `json:"count"`
}

func MapToBusinessConnectionResponse(bc *models.BusinessConnection) BusinessConnectionResponse {
	return BusinessConnectionResponse{
		ID:                   bc.ID,
		InitiatingBusinessID: bc.InitiatingBusinessID,
		ReceivingBusinessID:  bc.ReceivingBusinessID,
		ConnectionType:       bc.ConnectionType,
		Status:               bc.Status,
		InitiatedByUserID:    bc.InitiatedByUserID,
		Notes:                bc.Notes,
		CreatedAt:            bc.CreatedAt,
		UpdatedAt:            bc.UpdatedAt,
		InitiatingBusiness:   MapBusinessToResponse(&bc.InitiatingBusiness),
		ReceivingBusiness:    MapBusinessToResponse(&bc.ReceivingBusiness),
		InitiatedByUser:      MapUserToResponse(&bc.InitiatedByUser),
	}
}

func MapToBusinessConnectionsResponse(connections []models.BusinessConnection) BusinessConnectionsResponse {
	conns := make([]BusinessConnectionResponse, len(connections))
	for i, connection := range connections {
		conns[i] = MapToBusinessConnectionResponse(&connection)
	}

	return BusinessConnectionsResponse{
		Connections: conns,
		Count:       len(conns),
	}
}
