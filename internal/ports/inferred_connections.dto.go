package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type CreateInferredConnectionInput struct {
	SourceEntityType string  `json:"source_entity_type" validate:"required"`
	SourceEntityID   uint    `json:"source_entity_id" validate:"required"`
	TargetEntityType string  `json:"target_entity_type" validate:"required"`
	TargetEntityID   uint    `json:"target_entity_id" validate:"required"`
	ConnectionType   string  `json:"connection_type" validate:"required"`
	ConfidenceScore  float64 `json:"confidence_score" validate:"required"`
	ModelVersion     string  `json:"model_version"`
}
type InferredConnectionResponse struct {
	ID               uint      `json:"id"`
	SourceEntityType string    `json:"source_entity_type"`
	SourceEntityID   uint      `json:"source_entity_id"`
	TargetEntityType string    `json:"target_entity_type"`
	TargetEntityID   uint      `json:"target_entity_id"`
	ConnectionType   string    `json:"connection_type"`
	ConfidenceScore  float64   `json:"confidence_score"`
	ModelVersion     string    `json:"model_version"`
	CreatedAt        time.Time `json:"created_at"`
}
func MapInferredConnectionToResponse(ic *models.InferredConnection) InferredConnectionResponse {
	return InferredConnectionResponse{
		ID:               ic.ID,
		SourceEntityType: ic.SourceEntityType,
		SourceEntityID:   ic.SourceEntityID,
		TargetEntityType: ic.TargetEntityType,
		TargetEntityID:   ic.TargetEntityID,
		ConnectionType:   ic.ConnectionType,
		ConfidenceScore:  ic.ConfidenceScore,
		ModelVersion:     ic.ModelVersion,
		CreatedAt:        ic.CreatedAt,
	}
}
