package ports
import (
	"time"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)
type CreateBusinessTagInput struct {
	BusinessID  uint                   `json:"business_id" validate:"required"`
	TagType     models.BusinessTagType `json:"tag_type" validate:"required,oneof=client service specialty"`
	Description string                 `json:"description" validate:"required,min=1,max=100"`
}
type UpdateBusinessTagInput struct {
	TagType     *models.BusinessTagType `json:"tag_type,omitempty" validate:"omitempty,oneof=client service specialty"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,min=1,max=100"`
}
type BusinessTagResponse struct {
	ID          uint                   `json:"id"`
	BusinessID  uint                   `json:"business_id"`
	TagType     models.BusinessTagType `json:"tag_type"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	
	Business BusinessResponse `json:"business"`
}
type BusinessTagsResponse struct {
	Tags  []BusinessTagResponse `json:"tags"`
	Count int                   `json:"count"`
}
func MapToBusinessTagResponse(bt *models.BusinessTag) BusinessTagResponse {
	return BusinessTagResponse{
		ID:          bt.ID,
		BusinessID:  bt.BusinessID,
		TagType:     bt.TagType,
		Description: bt.Description,
		CreatedAt:   bt.CreatedAt,
		Business:    MapBusinessToResponse(&bt.Business),
	}
}
func MapToBusinessTagsResponse(businessTags []models.BusinessTag) BusinessTagsResponse {
	tags := make([]BusinessTagResponse, len(businessTags))
	for i, tag := range businessTags {
		tags[i] = MapToBusinessTagResponse(&tag)
	}
	return BusinessTagsResponse{
		Tags:  tags,
		Count: len(tags),
	}
}
